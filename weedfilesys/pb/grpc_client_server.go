package pb

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"math/rand/v2"
	"my_backend/weedfilesys/glog"
	"my_backend/weedfilesys/pb/filer_pb"
	"my_backend/weedfilesys/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Max_Message_Size = 1 << 30
)

var (
	grpcClients     = make(map[string]*versionedGrpcClient)
	grpcClientsLock sync.Mutex
)

type versionedGrpcClient struct {
	*grpc.ClientConn
	version  int
	errCount int
}

func WithFilerClient(streamingMode bool, signature int32, filer ServerAddress, option grpc.DialOption, fn func(client filer_pb.WeedfilesysFilerClient) error) error {
	return WithGrpcFilerClient(streamingMode, signature, filer, option, fn)
}

func WithGrpcFilerClient(streamingMode bool, signature int32, filerGrpcAddress ServerAddress, grpcDialOption grpc.DialOption, fn func(client filer_pb.WeedfilesysFilerClient) error) error {
	return WithGrpcClient(streamingMode, signature, func(conn *grpc.ClientConn) error {
		client := filer_pb.NewWeedfilesysFilerClient(conn)
		return fn(client)
	}, filerGrpcAddress.ToGrpcAddress(), false, grpcDialOption)
}

// 如果是在streamingMode，那么总是使用一个fresh的链接，否则就会复用同一个链接
func WithGrpcClient(streamingMode bool, signature int32, fn func(conn *grpc.ClientConn) error, address string, waitForReady bool, opts ...grpc.DialOption) error {
	if !streamingMode {
		vgc, err := getOrCreateConnection(address, waitForReady, opts...)
		if err != nil {
			return fmt.Errorf("getOrCreateConnection %s: %v", address, err)
		}
		executeErr := fn(vgc.ClientConn)
		if executeErr != nil {
			if strings.Contains(executeErr.Error(), "transport") ||
				strings.Contains(executeErr.Error(), "connection closed") {
				grpcClientsLock.Lock()
				if t, ok := grpcClients[address]; ok {
					if t.version == vgc.version {
						vgc.Close()
						delete(grpcClients, address)
					}
				}
				grpcClientsLock.Unlock()
			}
		}

		return executeErr
	} else {
		ctx := context.Background()
		if signature != 0 {
			// 👉 从一个键值对的 map 创建一个 MD（Metadata）对象，
			// 这个对象用于在 gRPC 中设置或传输元数据（如认证信息、用户信息、token 等）。
			md := metadata.New(map[string]string{"sw-client-id": fmt.Sprintf("%d", signature)})

			// 在原本的ctx的基础上创建一个带有出入站meta原数组的新的上下文
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
		grpcConnection, err := GrpcDial(ctx, address, waitForReady, opts...)
		if err != nil {
			return fmt.Errorf("fail to dial %s: %v", address, err)
		}
		defer grpcConnection.Close()
		executionErr := fn(grpcConnection)
		if executionErr != nil {
			return executionErr
		}
		return nil
	}
}

func getOrCreateConnection(address string, waitForReady bool, opts ...grpc.DialOption) (*versionedGrpcClient, error) {
	grpcClientsLock.Lock()
	defer grpcClientsLock.Unlock()

	existingConnection, found := grpcClients[address]
	if found {
		return existingConnection, nil
	}

	ctx := context.Background()
	grpcConnection, err := GrpcDial(ctx, address, waitForReady, opts...)
	if err != nil {
		return nil, fmt.Errorf("fail to dial %s, %v", address, err)
	}
	vgc := &versionedGrpcClient{
		ClientConn: grpcConnection,
		version:    rand.Int(),
		errCount:   0,
	}

	return vgc, nil
}

// 返回一个grpc客户端，waitForReady配置在链接中断或服务器不可达时进行RPC调用的行为
// 如果waitForReady 是false，且连接处于临时失败TRANSIENT_FAILURE，RPC会立即失败，否则RPC客户端会阻塞该调用，知道链接可以用
// 但是目前来看context没有用上，看看之后如果context有特定用处，就解决一下
func GrpcDial(ctx context.Context, address string, waitForReady bool, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var options []grpc.DialOption
	options = append(options, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(Max_Message_Size),
		grpc.MaxCallSendMsgSize(Max_Message_Size),
		grpc.WaitForReady(waitForReady)),

		// 返回一个链接设置用于指定链接存活的参数
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second, // 多久没有活动后客户端就会发送一个ping包
			Timeout:             20 * time.Second, //ping 之后等待pong响应的超时时间
			PermitWithoutStream: true,             // 是否允许在没有活跃RPC的情况下发送ping
		}),
	)
	for _, opt := range opts {
		if opt != nil {
			options = append(options, opt)
		}
	}

	return grpc.NewClient(address, options...)

}

func ServerToGrpcAddress(server string) (serverGrpcAddress string) {
	host, port, parseErr := hostAndPort(server)
	if parseErr != nil {
		glog.Fatalf("server address %s parse error: %v", server, parseErr)
	}
	grpcPort := int(port) + 10000
	return util.JoinHostPort(host, grpcPort)
}

func hostAndPort(address string) (host string, port uint64, err error) {
	colonIndex := strings.LastIndex(address, ":")
	if colonIndex < 0 {
		return "", 0, fmt.Errorf("server should have hostname:port format: %v", address)
	}
	port, err = strconv.ParseUint(address[colonIndex+1:], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("server port parse error: %w", err)
	}
	return address[:colonIndex], port, err
}
