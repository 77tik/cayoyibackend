package pb

import (
	"cayoyibackend/weedfilesys/glog"
	"reflect"
)

// ServerDiscovery 服务发现会找到至少一个服务实例
// 并且还会提供一些工具函数去刷新实例列表
// 🔹 SRV 记录格式
// SRV 记录的 DNS 查询名格式如下：
// _service._proto.name
// 其中：
//
// _service：服务名（比如 _filer, _sip, _ldap）
//
// _proto：协议类型（比如 _tcp, _udp）
//
// name：主域名（比如 mydomain.com）
//
// 例子：
// _filer._tcp.mydomain.com
type ServerDiscovery struct {
	list      []ServerAddress
	srvRecord *ServerSrvAddress
}

func NewServiceDiscoveryFromMap(m map[string]ServerAddress) (sd *ServerDiscovery) {
	sd = &ServerDiscovery{}
	for _, s := range m {
		sd.list = append(sd.list, s)
	}
	return sd
}

// DNS SRV：DNS服务记录
// 普通的 DNS 查询（A 记录、CNAME 记录）是：
// 给我 example.com，我想知道它的 IP 地址。
// SRV 记录则是：
// 给我在 example.com 域下提供某种服务（比如 filer 服务）的主机，它运行在哪个地址？端口是多少？优先级和负载情况如何？
// 可能会返回：
// 10 60 9333 filer1.mydomain.com.
// 20 40 9333 filer2.mydomain.com.
// 表示：
// filer1.mydomain.com:9333 提供 filer 服务，优先级为 10，权重为 60
//
// filer2.mydomain.com:9333 提供同样服务，优先级为 20，权重为 40
func (sd *ServerDiscovery) RefreshBySrvIfAvailable() {
	if sd.srvRecord == nil {
		return
	}

	// LookupSRV("", "", "_filer._tcp.mydomain.com")
	// 就能自动获取所有在该域下注册的 filer 服务节点，
	// 从而实现高可用的集群服务发现，而不需要硬编码 IP 地址或端口。
	newList, err := sd.srvRecord.LookUp()
	if err != nil {
		glog.V(0).Infof("failed to lookup SRV for %s: %v", *sd.srvRecord, err)
	}
	if newList == nil || len(newList) == 0 {
		glog.V(0).Infof("looked up SRV for %s, but found no well-formed names", *sd.srvRecord)
		return
	}
	// 不能用==，因为像map，struct之类的要用反射提供的这个函数来比较是否星等
	if !reflect.DeepEqual(sd.list, newList) {
		sd.list = newList
	}
}

// 获取全部实例，返回一个最新地址表的复制
func (sd *ServerDiscovery) GetInstances() (address []ServerAddress) {
	for _, s := range sd.list {
		address = append(address, s)
	}
	return address
}
func (sd *ServerDiscovery) GetInstancesAsStrings() (addresses []string) {
	for _, i := range sd.list {
		addresses = append(addresses, string(i))
	}
	return addresses
}
func (sd *ServerDiscovery) GetInstancesAsMap() (addresses map[string]ServerAddress) {
	addresses = make(map[string]ServerAddress)
	for _, i := range sd.list {
		addresses[string(i)] = i
	}
	return addresses
}
