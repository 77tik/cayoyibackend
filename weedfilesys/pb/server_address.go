package pb

import (
	"net"
	"strings"
)

type ServerAddress string
type ServerGrpcAddress string

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
