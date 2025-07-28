package backend

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/pb/master_pb"
	"cayoyibackend/weedfilesys/pb/volume_server_pb"
	"cayoyibackend/weedfilesys/util"
	"io"
	"os"
	"strings"
	"time"
)

// Weedfilesys 中后端存储抽象相关的实现，
// 主要负责管理和加载各种存储后端接口，
// 实现对不同存储后端（本地磁盘、云存储等）的统一访问。

// 后端抽象文件
type BackendStorageFile interface {
	// 支持读写
	io.ReaderAt
	io.WriterAt

	// 支持截断文件大小
	Truncate(off int64) error

	// 支持关闭
	io.Closer

	// 获取文件状态
	GetStat() (datSize int64, modTime time.Time, err error)
	Name() string

	// 同步磁盘
	Sync() error
}

// 抽象一个存储系统的能力，比如本地磁盘、S3、GCS 等
type BackendStorage interface {
	ToProperties() map[string]string                                                                                       //转成属性映射（方便配置序列化）
	NewStorageFile(key string, tierInfo *volume_server_pb.VolumeInfo) BackendStorageFile                                   //根据 key 和可选分层信息创建文件接口
	CopyFile(f *os.File, fn func(progressed int64, percentage float32) error) (key string, size int64, err error)          // 文件上传
	DownloadFile(fileName string, key string, fn func(progressed int64, percentage float32) error) (size int64, err error) // 文件下载
	DeleteFile(key string) (err error)                                                                                     // 删除文件
}

type StringProperties interface {
	GetString(key string) string
}
type StorageType string
type BackendStorageFactory interface {
	StorageType() StorageType
	BuildStorage(configuration StringProperties, configPrefix string, id string) (BackendStorage, error)
}

var (
	BackendStorageFactories = make(map[StorageType]BackendStorageFactory)
	BackendStorages         = make(map[string]BackendStorage)
)

// 加载配置初始化存储实例
// 从配置文件中读取 "storage.backend" 相关的配置节
// 遍历各存储类型和存储实例配置
// 跳过未启用的实例
// 对每个实例调用对应工厂构建存储实例
// 加入 BackendStorages 管理集合
// 若实例ID是 "default"，则也设置 BackendStorages 以类型名为 key 的快捷访问
func LoadConfiguration(config *util.ViperProxy) {

	StorageBackendPrefix := "storage.backend"

	for backendTypeName := range config.GetStringMap(StorageBackendPrefix) {
		backendStorageFactory, found := BackendStorageFactories[StorageType(backendTypeName)]
		if !found {
			glog.Fatalf("backend storage type %s not found", backendTypeName)
		}
		for backendStorageId := range config.GetStringMap(StorageBackendPrefix + "." + backendTypeName) {
			if !config.GetBool(StorageBackendPrefix + "." + backendTypeName + "." + backendStorageId + ".enabled") {
				continue
			}
			if _, found := BackendStorages[backendTypeName+"."+backendStorageId]; found {
				continue
			}
			backendStorage, buildErr := backendStorageFactory.BuildStorage(config,
				StorageBackendPrefix+"."+backendTypeName+"."+backendStorageId+".", backendStorageId)
			if buildErr != nil {
				glog.Fatalf("fail to create backend storage %s.%s", backendTypeName, backendStorageId)
			}
			BackendStorages[backendTypeName+"."+backendStorageId] = backendStorage
			if backendStorageId == "default" {
				BackendStorages[backendTypeName] = backendStorage
			}
		}
	}

}

// 用于Volume Server 从Master接收到的远程存储配置来初始化
func LoadFromPbStorageBackends(storageBackends []*master_pb.StorageBackend) {

	for _, storageBackend := range storageBackends {
		backendStorageFactory, found := BackendStorageFactories[StorageType(storageBackend.Type)]
		if !found {
			glog.Warningf("storage type %s not found", storageBackend.Type)
			continue
		}
		if _, found := BackendStorages[storageBackend.Type+"."+storageBackend.Id]; found {
			continue
		}
		backendStorage, buildErr := backendStorageFactory.BuildStorage(newProperties(storageBackend.Properties), "", storageBackend.Id)
		if buildErr != nil {
			glog.Fatalf("fail to create backend storage %s.%s", storageBackend.Type, storageBackend.Id)
		}
		BackendStorages[storageBackend.Type+"."+storageBackend.Id] = backendStorage
		if storageBackend.Id == "default" {
			BackendStorages[storageBackend.Type] = backendStorage
		}
	}
}

// 配置属性的封装
type Properties struct {
	m map[string]string
}

func newProperties(m map[string]string) *Properties {
	return &Properties{m: m}
}

func (p *Properties) GetString(key string) string {
	if v, found := p.m[key]; found {
		return v
	}
	return ""
}

// 导出当前已加载的存储实例为 protobuf 配置列表
// 遍历所有 BackendStorages
// 解析 key（类型和实例ID）
// 调用存储实例的 ToProperties 获取配置属性
// 封装成 protobuf 对象返回
func ToPbStorageBackends() (backends []*master_pb.StorageBackend) {
	for sName, s := range BackendStorages {
		sType, sId := BackendNameToTypeId(sName)
		if sType == "" {
			continue
		}
		backends = append(backends, &master_pb.StorageBackend{
			Type:       sType,
			Id:         sId,
			Properties: s.ToProperties(),
		})
	}
	return
}

// 把存储实例名称解析成存储类型和实例ID
func BackendNameToTypeId(backendName string) (backendType, backendId string) {
	parts := strings.Split(backendName, ".")
	if len(parts) == 1 {
		return backendName, "default"
	}
	if len(parts) != 2 {
		return
	}

	backendType, backendId = parts[0], parts[1]
	return
}
