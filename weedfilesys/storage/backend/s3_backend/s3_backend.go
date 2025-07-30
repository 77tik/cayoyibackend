package s3_backend

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/pb/volume_server_pb"
	"cayoyibackend/weedfilesys/storage/backend"
	"cayoyibackend/weedfilesys/util"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/uuid"
	"io"
	"os"
	"strings"
	"time"
)

// “S3 存储”是指基于 Amazon S3 协议（Simple Storage Service） 的对象存储服务。
// 它是一种 云端的、可扩展的、通过 HTTP 接口访问的文件存储方式，被广泛用于存储图片、视频、日志、备份文件、大数据等不可变的静态文件。
// | 概念          | 含义                                                    |
// | ------------ | ----------------------------------------------------- |
// | **Bucket**   | 存储桶，相当于一个顶级的文件夹（项目/应用的命名空间）                           |
// | **Object**   | 存储的具体内容，例如 `.jpg`、`.txt`、`.dat` 等文件                   |
// | **Key**      | 对象的路径标识符，相当于文件名（可带路径如 `folder/file.txt`）              |
// | **Endpoint** | S3 服务的地址（AWS 是 `s3.amazonaws.com`，MinIO/COS/OSS 是自定义） |

func init() {
	//注册 S3 存储类型的构造器：
	backend.BackendStorageFactories["s3"] = &S3BackendFactory{}
}

type S3BackendFactory struct {
}

func (factory *S3BackendFactory) StorageType() backend.StorageType {
	return backend.StorageType("s3")
}
func (factory *S3BackendFactory) BuildStorage(configuration backend.StringProperties, configPrefix string, id string) (backend.BackendStorage, error) {
	return newS3BackendStorage(configuration, configPrefix, id)
}

type S3BackendStorage struct {
	id                    string
	aws_access_key_id     string
	aws_secret_access_key string
	region                string
	bucket                string
	endpoint              string
	storageClass          string
	forcePathStyle        bool
	conn                  s3iface.S3API
}

// 从配置中读取 AWS Key/Secret、Bucket、Region、Endpoint（支持自定义 MinIO 之类）
// 创建 AWS S3 连接 s3iface.S3API
// 可通过 force_path_style = true 支持路径风格 URL（适配 MinIO）
func newS3BackendStorage(configuration backend.StringProperties, configPrefix string, id string) (s *S3BackendStorage, err error) {
	s = &S3BackendStorage{}
	s.id = id
	s.aws_access_key_id = configuration.GetString(configPrefix + "aws_access_key_id")
	s.aws_secret_access_key = configuration.GetString(configPrefix + "aws_secret_access_key")
	s.region = configuration.GetString(configPrefix + "region")
	s.bucket = configuration.GetString(configPrefix + "bucket")
	s.endpoint = configuration.GetString(configPrefix + "endpoint")
	s.storageClass = configuration.GetString(configPrefix + "storage_class")
	s.forcePathStyle = util.ParseBool(configuration.GetString(configPrefix+"force_path_style"), true)
	if s.storageClass == "" {
		s.storageClass = "STANDARD_IA"
	}

	s.conn, err = createSession(s.aws_access_key_id, s.aws_secret_access_key, s.region, s.endpoint, s.forcePathStyle)

	glog.V(0).Infof("created backend storage s3.%s for region %s bucket %s", s.id, s.region, s.bucket)
	return
}

func (s *S3BackendStorage) ToProperties() map[string]string {
	m := make(map[string]string)
	m["aws_access_key_id"] = s.aws_access_key_id
	m["aws_secret_access_key"] = s.aws_secret_access_key
	m["region"] = s.region
	m["bucket"] = s.bucket
	m["endpoint"] = s.endpoint
	m["storage_class"] = s.storageClass
	m["force_path_style"] = util.BoolToString(s.forcePathStyle)
	return m
}

func (s *S3BackendStorage) NewStorageFile(key string, tierInfo *volume_server_pb.VolumeInfo) backend.BackendStorageFile {
	if strings.HasPrefix(key, "/") {
		key = key[1:]
	}

	f := &S3BackendStorageFile{
		backendStorage: s,
		key:            key,
		tierInfo:       tierInfo,
	}

	return f
}

// 生成 UUID 作为对象 key
// 上传本地 .dat 文件到 S3，对应为一个对象
// 支持进度回调函数 fn(...) 报告进度
func (s *S3BackendStorage) CopyFile(f *os.File, fn func(progressed int64, percentage float32) error) (key string, size int64, err error) {
	randomUuid, _ := uuid.NewRandom()
	key = randomUuid.String()

	glog.V(1).Infof("copying dat file of %s to remote s3.%s as %s", f.Name(), s.id, key)

	util.Retry("upload to S3", func() error {
		size, err = uploadToS3(s.conn, f.Name(), s.bucket, key, s.storageClass, fn)
		return err
	})

	return
}

// 将远程对象下载为本地文件，支持回调进度
func (s *S3BackendStorage) DownloadFile(fileName string, key string, fn func(progressed int64, percentage float32) error) (size int64, err error) {

	glog.V(1).Infof("download dat file of %s from remote s3.%s as %s", fileName, s.id, key)

	size, err = downloadFromS3(s.conn, fileName, s.bucket, key, fn)

	return
}

// 删除 S3 上的对应对象（volume 数据文件）
func (s *S3BackendStorage) DeleteFile(key string) (err error) {

	glog.V(1).Infof("delete dat file %s from remote", key)

	err = deleteFromS3(s.conn, s.bucket, key)

	return
}

// 将s3实现为内部的调用需要的backend接口
type S3BackendStorageFile struct {
	backendStorage *S3BackendStorage
	key            string
	tierInfo       *volume_server_pb.VolumeInfo
}

func (s3backendStorageFile S3BackendStorageFile) ReadAt(p []byte, off int64) (n int, err error) {
	datSize, _, _ := s3backendStorageFile.GetStat()

	if datSize > 0 && off >= datSize {
		return 0, io.EOF
	}

	bytesRange := fmt.Sprintf("bytes=%d-%d", off, off+int64(len(p))-1)

	getObjectOutput, getObjectErr := s3backendStorageFile.backendStorage.conn.GetObject(&s3.GetObjectInput{
		Bucket: &s3backendStorageFile.backendStorage.bucket,
		Key:    &s3backendStorageFile.key,
		Range:  &bytesRange,
	})

	if getObjectErr != nil {
		return 0, fmt.Errorf("bucket %s GetObject %s: %v", s3backendStorageFile.backendStorage.bucket, s3backendStorageFile.key, getObjectErr)
	}
	defer getObjectOutput.Body.Close()

	// glog.V(3).Infof("read %s %s", s3backendStorageFile.key, bytesRange)
	// glog.V(3).Infof("content range: %s, contentLength: %d", *getObjectOutput.ContentRange, *getObjectOutput.ContentLength)

	var readCount int
	for {
		p = p[readCount:]
		readCount, err = getObjectOutput.Body.Read(p)
		n += readCount

		if err != nil {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func (s3backendStorageFile S3BackendStorageFile) WriteAt(p []byte, off int64) (n int, err error) {
	panic("not implemented")
}

func (s3backendStorageFile S3BackendStorageFile) Truncate(off int64) error {
	panic("not implemented")
}

func (s3backendStorageFile S3BackendStorageFile) Close() error {
	return nil
}

func (s3backendStorageFile S3BackendStorageFile) GetStat() (datSize int64, modTime time.Time, err error) {

	files := s3backendStorageFile.tierInfo.GetFiles()

	if len(files) == 0 {
		err = fmt.Errorf("remote file info not found")
		return
	}

	datSize = int64(files[0].FileSize)
	modTime = time.Unix(int64(files[0].ModifiedTime), 0)

	return
}

func (s3backendStorageFile S3BackendStorageFile) Name() string {
	return s3backendStorageFile.key
}

func (s3backendStorageFile S3BackendStorageFile) Sync() error {
	return nil
}
