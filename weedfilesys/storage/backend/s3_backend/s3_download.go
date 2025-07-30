package s3_backend

import (
	"cayoyibackend/weedfilesys/glog"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"sync/atomic"
)

// 从 S3 指定 Bucket + Key 下载文件到本地 destFileName
// 提供下载进度通知（fn 是一个回调函数）
func downloadFromS3(sess s3iface.S3API, destFileName string, sourceBucket string, sourceKey string,
	fn func(progressed int64, percentage float32) error) (fileSize int64, err error) {

	// 调用 getFileSize() 获取远程 S3 文件大小。
	// 打开本地目标文件用于写入。
	// 创建一个带进度感知的 Writer（s3DownloadProgressedWriter）
	// 使用 AWS SDK 提供的 s3manager.Downloader 进行并发分片下载。
	// 每次写入触发一次 fn() 回调，用于显示进度条或记录日志。
	fileSize, err = getFileSize(sess, sourceBucket, sourceKey)
	if err != nil {
		return
	}

	//open the file
	f, err := os.OpenFile(destFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %q, %v", destFileName, err)
	}
	defer f.Close()

	// Create a downloader with the session and custom options
	downloader := s3manager.NewDownloaderWithClient(sess, func(u *s3manager.Downloader) {
		u.PartSize = int64(64 * 1024 * 1024)
		u.Concurrency = 5
	})

	fileWriter := &s3DownloadProgressedWriter{
		fp:      f,
		size:    fileSize,
		written: 0,
		fn:      fn,
	}

	// Download the file from S3.
	fileSize, err = downloader.Download(fileWriter, &s3.GetObjectInput{
		Bucket: aws.String(sourceBucket),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		return fileSize, fmt.Errorf("failed to download /buckets/%s%s to %s: %v", sourceBucket, sourceKey, destFileName, err)
	}

	glog.V(1).Infof("downloaded file %s\n", destFileName)

	return
}

// adapted from https://github.com/aws/aws-sdk-go/pull/1868
// and https://petersouter.xyz/s3-download-progress-bar-in-golang/
// 这是一个 包装 os.File 的结构体，实现了 WriteAt() 方法，并在每次写入时：
// 原子性地累加写入字节数 w.written
// 调用进度回调 fn(written, written/size)
// 目的是为了在下载中实时获取进度。
type s3DownloadProgressedWriter struct {
	size    int64
	written int64
	fn      func(progressed int64, percentage float32) error
	fp      *os.File
}

// ✅ 为什么要实现 WriteAt 接口？
// 为了兼容 AWS SDK 的 s3manager.Downloader。
// 🔧 s3manager.Downloader 要求目标对象实现 io.WriterAt
// s3manager.Downloader 会把远程 S3 文件切成 N 个“块”，每个块开一个 goroutine 来下载。
//
// 下载时，它会调用你传入的 WriterAt.WriteAt(p []byte, offset int64)，写入对应位置。
// 传入了 s3DownloadProgressedWriter，它底层使用了 os.File.WriteAt() 写入到本地文件，并加了进度统计逻辑。
func (w *s3DownloadProgressedWriter) WriteAt(p []byte, off int64) (int, error) {
	n, err := w.fp.WriteAt(p, off)
	if err != nil {
		return n, err
	}

	// Got the length have read( or means has uploaded), and you can construct your message
	atomic.AddInt64(&w.written, int64(n))

	if w.fn != nil {
		written := w.written
		if err := w.fn(written, float32(written*100)/float32(w.size)); err != nil {
			return n, err
		}
	}

	return n, err
}

// 通过 s3.HeadObject() 获取远程文件的 ContentLength，用于：
// 预估总大小
// 计算进度百分比
func getFileSize(svc s3iface.S3API, bucket string, key string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}
