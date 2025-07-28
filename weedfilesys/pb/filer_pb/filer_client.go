package filer_pb

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/util"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	OS_UID = uint32(os.Getuid())
	OS_GID = uint32(os.Getgid())
)

type FilerClient interface {
	WithFilerClient(streamingMode bool, fn func(client WeedfilesysFilerClient) error) error
	AdjustedUrl(location *Location) string
	GetDataCenter() string
}

func GetEntry(ctx context.Context, filerClient FilerClient, fullFilePath util.FullPath) (entry *Entry, err error) {
	dir, name := fullFilePath.DirAndName()
	err = filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {
		requst := &LookupDirectoryEntryRequest{
			Directory: dir,
			Name:      name,
		}

		resp, err := LookupEntry(ctx, client, requst)
		if err != nil {
			glog.V(3).InfofCtx(ctx, "read %s %v: %v", fullFilePath, resp, err)
			return err
		}

		if resp.Entry == nil {
			return nil
		}
		entry = resp.Entry
		return nil
	})

	return
}

type EachEntryFunction func(entry *Entry, isLast bool) error

// Entry 目录项
// 当我们列出一个目录时，例如：
// /home/user/docs
// ├── a.txt
// ├── b.txt
// └── subdir/
// 这个目录下的内容就是 3 个 entries（条目）：
//
// 一个是文件 a.txt
//
// 一个是文件 b.txt
//
// 一个是子目录 subdir
// 还添加了分页控制，处理了每一页最多取多少
func ReadDirAllEntries(ctx context.Context, filerClient FilerClient, fullDirPath util.FullPath, prefix string, fn EachEntryFunction) (err error) {
	var counter uint32
	var startFrom string
	var counterFunc = func(entry *Entry, isLast bool) error {
		counter++
		startFrom = entry.Name
		return fn(entry, isLast)
	}

	// 最多读取10000个entry
	var paginationLimit uint32 = 10000
	if err = doList(ctx, filerClient, fullDirPath, prefix, counterFunc, "", false, paginationLimit); err != nil {
		return err
	}

	// 如果读取的个数正好是10000个，那么可能有更多的entry，于是继续下一轮分页读取，否则读取完成后返回
	for counter == paginationLimit {
		counter = 0
		if err = doList(ctx, filerClient, fullDirPath, prefix, counterFunc, startFrom, false, paginationLimit); err != nil {
			return err
		}
	}
	return nil
}

func doList(ctx context.Context, filerClient FilerClient, fullDirPath util.FullPath, prefix string, fn EachEntryFunction, startFrom string, inclusive bool, limit uint32) (err error) {
	return filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {
		return doWeedfilesysList(ctx, client, fullDirPath, prefix, fn, startFrom, inclusive, limit)
	})
}

func List(ctx context.Context, filerClient FilerClient, parentDirectoryPath, prefix string, fn EachEntryFunction, startFrom string, inclusive bool, limit uint32) (err error) {
	return filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {
		return doWeedfilesysList(ctx, client, util.FullPath(parentDirectoryPath), prefix, fn, startFrom, inclusive, limit)
	})
}

func doWeedfilesysList(ctx context.Context, client WeedfilesysFilerClient, fullDirPath util.FullPath, prefix string, fn EachEntryFunction, startFrom string, inclusive bool, limit uint32) (err error) {
	// 冗余限制，用于正确判断是否为最后一个文件
	// 比如你在处理分块上传文件，如果你要判断：“我是不是已经收到了最后一个文件块”，
	// 但由于网络抖动、延迟、乱序等问题，你可能会设置一个“冗余限制”，等超过某个阈值后才真正确认是“最后一块”。
	redLimit := limit
	// 当我们设置limit=10的时候，想分页读取目录的10个entry，但如果正好目录有超过10个entry，那么我们需要一种方法判断
	// 做法就是多请求一个(limit+1),如果返回了11个entry就代表后面还有
	if limit < math.MaxInt32 && limit != 0 {
		redLimit = limit + 1
	}
	if redLimit > math.MaxInt32 {
		redLimit = math.MaxInt32
	}

	request := &ListEntriesRequest{
		Directory:          string(fullDirPath),
		Prefix:             prefix,
		StartFromFileName:  startFrom,
		Limit:              redLimit,
		InclusiveStartFrom: inclusive,
	}

	glog.V(4).InfofCtx(ctx, "doList request: %v", request)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := client.ListEntries(ctx, request)
	if err != nil {
		return fmt.Errorf("list %s: %v", fullDirPath, err)
	}

	// 请求 limit+1 条 → 多拿一条来判断是否还有下一页
	//    ↓
	//for 每条 entry:
	//    - 如果是最后一条（第 limit+1 个）→ 不调用 fn
	//    - 如果是倒数第二条（第 limit 个）→ 调用 fn(prev, false)
	//    - 如果返回的 entry 不满 limit → 最后一条调用 fn(prev, true)
	var prevEntry *Entry
	count := 0
	for {
		resp, recvErr := stream.Recv()
		if recvErr != nil {
			if recvErr == io.EOF {
				if prevEntry != nil {
					if ee := fn(prevEntry, true); ee != nil {
						return ee
					}
				}
				break
			} else {
				return recvErr
			}
		}
		if prevEntry != nil {
			if ee := fn(prevEntry, false); ee != nil {
				return ee
			}
		}
		prevEntry = resp.Entry
		count++

		// 超出limit+1，丢弃最后一个冗余
		if count > int(limit) && limit != 0 {
			prevEntry = nil
		}

	}
	return nil
}

func WeedfilesysList(ctx context.Context, client WeedfilesysFilerClient, parentDirectoryPath, prefix string, fn EachEntryFunction, startFrom string, inclusive bool, limit uint32) (err error) {
	return doWeedfilesysList(ctx, client, util.FullPath(parentDirectoryPath), prefix, fn, startFrom, inclusive, limit)
}

func Exists(ctx context.Context, filerClient FilerClient, parentDirectoryPath string, entryName string, isDirectory bool) (exists bool, err error) {
	err = filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {
		request := &LookupDirectoryEntryRequest{
			Directory: parentDirectoryPath,
			Name:      entryName,
		}

		glog.V(4).InfofCtx(ctx, "exists entry %v/%v: %v", parentDirectoryPath, entryName, request)
		resp, err := LookupEntry(ctx, client, request)
		if err != nil {
			if err == ErrNotFound {
				exists = false
				return nil
			}
			glog.V(0).InfofCtx(ctx, "exists entry %v: %v", request, err)
			return fmt.Errorf("exists entry %s/%s: %v", parentDirectoryPath, entryName, err)
		}

		exists = resp.Entry.IsDirectory == isDirectory

		return nil
	})

	return
}

func Touch(ctx context.Context, filerClient FilerClient, parentDirectoryPath string, entryName string, entry *Entry) (err error) {
	return filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {
		request := &UpdateEntryRequest{
			Directory: parentDirectoryPath,
			Entry:     entry,
		}
		glog.V(4).InfofCtx(ctx, "touch entry %v/%v: %v", parentDirectoryPath, entryName, request)
		if err := UpdateEntry(ctx, client, request); err != nil {
			glog.V(0).InfofCtx(ctx, "touch exists entry %v: %v", request, err)
			return fmt.Errorf("touch exists entry %s/%s: %v", parentDirectoryPath, entryName, err)
		}

		return nil
	})
}

func Mkdir(ctx context.Context, filerClient FilerClient, parentDirectoryPath string, dirName string, fn func(entry *Entry)) (err error) {
	return filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {
		return DoMkdir(ctx, client, parentDirectoryPath, dirName, fn)
	})
}
func DoMkdir(ctx context.Context, client WeedfilesysFilerClient, parentDirectoryPath string, dirName string, fn func(entry *Entry)) (err error) {
	entry := &Entry{
		Name:        dirName,
		IsDirectory: true,
		Attributes: &FuseAttributes{
			Mime:     strconv.FormatInt(time.Now().Unix(), 10),
			Crtime:   time.Now().Unix(),
			FileMode: uint32(os.ModeDir | 0777),
			Uid:      OS_UID,
			Gid:      OS_GID,
		},
	}

	if fn != nil {
		fn(entry)
	}

	request := &CreateEntryRequest{
		Directory: parentDirectoryPath,
		Entry:     entry,
	}

	glog.V(1).InfofCtx(ctx, "mkdir: %v", request)
	if err := CreateEntry(ctx, client, request); err != nil {
		glog.V(0).InfofCtx(ctx, "mkdir %v: %v", request, err)
		return fmt.Errorf("mkdir %s/%s: %v", parentDirectoryPath, dirName, err)
	}

	return nil
}

func MkFile(ctx context.Context, filerClient FilerClient, parentDirectoryPath string, fileName string, chunks []*FileChunk, fn func(entry *Entry)) error {
	return filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {

		entry := &Entry{
			Name:        fileName,
			IsDirectory: false,
			Attributes: &FuseAttributes{
				Mtime:    time.Now().Unix(),
				Crtime:   time.Now().Unix(),
				FileMode: uint32(0770),
				Uid:      OS_UID,
				Gid:      OS_GID,
			},
			Chunks: chunks,
		}

		if fn != nil {
			fn(entry)
		}

		request := &CreateEntryRequest{
			Directory: parentDirectoryPath,
			Entry:     entry,
		}

		glog.V(1).InfofCtx(ctx, "create file: %s/%s", parentDirectoryPath, fileName)
		if err := CreateEntry(ctx, client, request); err != nil {
			glog.V(0).InfofCtx(ctx, "create file %v:%v", request, err)
			return fmt.Errorf("create file %s/%s: %v", parentDirectoryPath, fileName, err)
		}

		return nil
	})
}

func Remove(ctx context.Context, filerClient FilerClient, parentDirectoryPath, name string, isDeleteData, isRecursive, ignoreRecursiveErr, isFromOtherCluster bool, signatures []int32) error {
	return filerClient.WithFilerClient(false, func(client WeedfilesysFilerClient) error {
		return DoRemove(ctx, client, parentDirectoryPath, name, isDeleteData, isRecursive, ignoreRecursiveErr, isFromOtherCluster, signatures)
	})
}

func DoRemove(ctx context.Context, client WeedfilesysFilerClient, parentDirectoryPath string, name string, isDeleteData bool, isRecursive bool, ignoreRecursiveErr bool, isFromOtherCluster bool, signatures []int32) error {
	deleteEntryRequest := &DeleteEntryRequest{
		Directory:            parentDirectoryPath,
		Name:                 name,
		IsDeleteData:         isDeleteData,
		IsRecursive:          isRecursive,
		IgnoreRecursiveError: ignoreRecursiveErr,
		IsFromOtherCluster:   isFromOtherCluster,
		Signatures:           signatures,
	}
	if resp, err := client.DeleteEntry(ctx, deleteEntryRequest); err != nil {
		if strings.Contains(err.Error(), ErrNotFound.Error()) {
			return nil
		}
		return err
	} else {
		if resp.Error != "" {
			if strings.Contains(resp.Error, ErrNotFound.Error()) {
				return nil
			}
			return errors.New(resp.Error)
		}
	}

	return nil
}
