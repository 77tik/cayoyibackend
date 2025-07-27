package pb

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ServerAddress string
type ServerGrpcAddress string
type ServerSrvAddress string

// 将ServerAddress 解析为标准的Grpc地址格式：
// 在一些分布式系统中，服务注册地址可能带有端口段信息，比如：
// 127.0.0.1:8888.1 表示第一个实例占用了端口段
// 假设有一台服务器，想在上面部署10个grpc服务实例，传统做法是为每个实例分配一个完整的端口：
//
//	127.0.0.1:8001, 127.0.0.1:8002, 127.0.0.1:8003 ....
//
// 端口段是一个伪端口，需要转换成真正的端口，这就是此函数的用法
// "127.0.0.1:8000.1" → "127.0.0.1:1"
func (sa ServerAddress) ToGrpcAddress() string {
	portsSepIndex := strings.LastIndex(string(sa), ":")
	if portsSepIndex < 0 {
		return string(sa)
	}
	if portsSepIndex+1 >= len(sa) {
		return string(sa)
	}
	ports := string(sa[portsSepIndex+1:])
	sepIndex := strings.LastIndex(ports, ".")
	if sepIndex >= 0 {
		host := string(sa[0:portsSepIndex])
		return net.JoinHostPort(host, ports[sepIndex+1:])
	}
	return ServerToGrpcAddress(string(sa))
}

// “SRV” 是 Service Record（服务记录） 的缩写，
// 是 DNS（域名系统）中的一种特殊记录类型，用于告诉客户端：
// 在某个域名下，提供特定服务的主机是哪些，以及它们的端口号和优先级。
func (r ServerSrvAddress) LookUp() (addresses []ServerAddress, err error) {
	_, records, lookupErr := net.LookupSRV("", "", string(r))
	if lookupErr != nil {
		err = fmt.Errorf("lookup SRV address %s: %v", r, lookupErr)
	}
	for _, srv := range records {
		address := fmt.Sprintf("%s:%d", srv.Target, int(srv.Port))
		addresses = append(addresses, ServerAddress(address))
	}
	return
}
func NewServerAddressWithGrpcPort(address string, grpcPort int) ServerAddress {
	if grpcPort == 0 {
		return ServerAddress(address)
	}
	_, port, _ := hostAndPort(address)
	if uint64(grpcPort) == port+10000 {
		return ServerAddress(address)
	}
	return ServerAddress(address + "." + strconv.Itoa(grpcPort))
}
