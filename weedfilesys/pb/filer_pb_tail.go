package pb

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"my_backend/weedfilesys/glog"
	"my_backend/weedfilesys/pb/filer_pb"
	"my_backend/weedfilesys/util"
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
	DirectoriesToWatch     []string
	StartTsNs              int64
	StopTsNs               int64
	EventErrorType         EventErrorType
}

// 处理订阅元信息返回值的回调函数
type ProcessMetadataFunc func(resp *filer_pb.SubscribeMetadataResponse) error

func FollowMetadata(filerServer ServerAddress, grpcDialOption grpc.DialOption, option *MetadataFollowOption, processEventFn ProcessMetadataFunc) error {
	err := WithFilerClient(true, option.SelfSignature, filerServer, grpcDialOption, makeSubscribeMetadataFunc(option, processEventFn))
	if err != nil {
		return fmt.Errorf("subscribing filer meta change: %w", err)
	}
	return nil
}

func WithFilerClientFollowMetadata(filerClient filer_pb.FilerClient, option *MetadataFollowOption, processEventFn ProcessMetadataFunc) error {

	err := filerClient.WithFilerClient(true, makeSubscribeMetadataFunc(option, processEventFn))
	if err != nil {
		return fmt.Errorf("subscribing filer meta change: %w", err)
	}

	return nil
}

// 接收原信息设置和处理原信息的回调函数
func makeSubscribeMetadataFunc(option *MetadataFollowOption, processEventFn ProcessMetadataFunc) func(client filer_pb.WeedfilesysFilerClient) error {
	return func(client filer_pb.WeedfilesysFilerClient) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stream, err := client.SubscribeMetadata(ctx, &filer_pb.SubscribeMetadataRequest{
			ClientName:   option.ClientName,
			PathPrefix:   option.PathPrefix,
			PathPrefixes: option.AdditionalPathPrefixes,
			Directories:  option.DirectoriesToWatch,
			SinceNs:      option.StartTsNs,
			Signature:    option.SelfSignature,
			ClientId:     option.ClientId,
			ClientEpoch:  option.ClientEpoch,
			UntilNs:      option.StopTsNs,
		})

		if err != nil {
			return fmt.Errorf("subscribe: %w", err)
		}

		for {
			resp, listenErr := stream.Recv()
			if listenErr == io.EOF {
				return nil
			}
			if listenErr != nil {
				return listenErr
			}

			if err := processEventFn(resp); err != nil {
				switch option.EventErrorType {
				case TrivialOnError:
					glog.Errorf("process %v: %v", resp, err)
				case FatalOnError:
					glog.Fatalf("process %v: %v", resp, err)
				case RetryForeverOnError:
					util.RetryUntil("followMetaUpdates", func() error { return processEventFn(resp) }, func(err error) bool {
						glog.Errorf("process %v: %v", resp, err)
						return true
					})
				case DontLogError:
				//pass
				default:
					glog.Errorf("process %v: %v", resp, err)
				}

			}
			option.StartTsNs = resp.TsNs
		}
	}
}
