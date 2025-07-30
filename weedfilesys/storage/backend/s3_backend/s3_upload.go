package s3_backend

import (
	"cayoyibackend/weedfilesys/glog"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"sync"
)

// 上传指定文件到 S3（支持分片、多并发）
func uploadToS3(sess s3iface.S3API, filename string, destBucket string, destKey string, storageClass string, fn func(progressed int64, percentage float32) error) (fileSize int64, err error) {

	//打开并读取文件大小
	f, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %q, %v", filename, err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to stat file %q, %v", filename, err)
	}
	fileSize = info.Size()

	// 动态调整分片大小
	// AWS S3 要求每个分片最小为 5MB，最大 5GB，最多 10,000 个分片。这个逻辑是为了控制并发上传的分片数不过多（最多 1000 个左右）。
	partSize := int64(64 * 1024 * 1024) // The minimum/default allowed part size is 5MB
	for partSize*1000 < fileSize {
		partSize *= 4
	}

	// 使用 AWS SDK 的分片上传器
	// 设置并发上传分片数为 5
	uploader := s3manager.NewUploaderWithClient(sess, func(u *s3manager.Uploader) {
		u.PartSize = partSize
		u.Concurrency = 5
	})

	// 将 os.File 封装成自定义 reader，在每次读取分片时能统计上传进度，并调用 fn(progressed, percentage) 函数。
	fileReader := &s3UploadProgressedReader{
		fp:      f,
		size:    fileSize,
		signMap: map[int64]struct{}{},
		fn:      fn,
	}

	// 实际发起上传操作，失败时抛错，成功后日志中输出上传位置（result.Location）。
	var result *s3manager.UploadOutput
	result, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:       aws.String(destBucket),
		Key:          aws.String(destKey),
		Body:         fileReader,
		StorageClass: aws.String(storageClass),
	})

	//in case it fails to upload
	if err != nil {
		return 0, fmt.Errorf("failed to upload file %s: %v", filename, err)
	}
	glog.V(1).Infof("file %s uploaded to %s\n", filename, result.Location)

	return
}

// adapted from https://github.com/aws/aws-sdk-go/pull/1868
// https://github.com/aws/aws-sdk-go/blob/main/example/service/s3/putObjectWithProcess/putObjWithProcess.go
// 自定义上传进度Reader
type s3UploadProgressedReader struct {
	fp      *os.File           //目标文件句柄
	size    int64              //文件总大小
	read    int64              //已上传字节数
	signMap map[int64]struct{} //用于避免重复记录偏移
	mux     sync.Mutex
	fn      func(progressed int64, percentage float32) error //自定义进度回调函数
}

//🔍 为什么要实现它们？
//因为 AWS S3 的 SDK（s3manager.Uploader）在上传文件时，会用到这些接口 来：
//
//进行文件 MD5 校验签名时（调用 ReadAt）
//
//正式上传数据时（调用 ReadAt）
//
//需要重试或跳转位置时（调用 Seek）
//
//文件过小时直接串行读取（调用 Read）

// 连续读取
func (r *s3UploadProgressedReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

// 按偏移量读取
func (r *s3UploadProgressedReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	r.mux.Lock()
	// Ignore the first signature call
	if _, ok := r.signMap[off]; ok {
		r.read += int64(n)
	} else {
		r.signMap[off] = struct{}{}
	}
	r.mux.Unlock()

	if r.fn != nil {
		read := r.read
		if err := r.fn(read, float32(read*100)/float32(r.size)); err != nil {
			return n, err
		}
	}

	return n, err
}

// 移动读写指针，用于定位位置
func (r *s3UploadProgressedReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}
