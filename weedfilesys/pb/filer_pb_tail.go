package pb

import (
	"my_backend/weedfilesys/pb/filer_pb"
)

// 这个 Go 代码文件定义了在 weedfilesys 系统中订阅和处理元数据变更事件的客户端逻辑，
// 核心是通过 gRPC 连接到 Filer 服务，调用它的 SubscribeMetadata 流式 RPC，持续接收元数据更新，
// 并用用户提供的函数处理这些事件。

type EventErrorType int

// 用于控制在处理元数据出错时的策略
const (
	TrivialOnError      EventErrorType = iota //记录错误，继续运行
	FatalOnError                              //出错直接fatal（终止程序）
	RetryForeverOnError                       // 遇到错误一直重试
	DontLogError                              // 出错不记录日志
)

// 订阅结构体
type MetadataFollowOption struct {
	ClientName             string
	ClientId               int32
	ClientEpoch            int32
	SelfSignature          int32
	PathPrefix             string
	AdditionalPathPrefixes []string
}

type ProcessMetadataFunc func(resp *filer_pb.SubscribeMetadataResponse) error

//func FollowMetadata(filerServer ServerAddress, grpcDialOption grpc.DialOption, option *MetadataFollowOption, processEventFn ProcessMetadataFunc) error {
//
//}
