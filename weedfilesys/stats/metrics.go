package stats

import (
	"cayoyibackend/weedfilesys/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Prometheus 是用来“采集、存储、查询和可视化系统运行状态指标（Metrics）”的工具。
//
//📦 Prometheus 能做什么？
//功能	描述
//⏱️ 指标采集	自动定期向目标服务发起 HTTP 请求，采集它暴露的指标（如：CPU使用率、请求次数等）
//💾 数据存储	自带时间序列数据库，存储每个指标的历史值（带时间戳）
//🔍 指标查询	提供强大的 PromQL 查询语言，灵活筛选和聚合指标
//📊 可视化	可以配合 Grafana 展示漂亮的图表，也支持自带的简易 Web UI
//🚨 告警系统	支持设置告警规则，当某个指标满足特定条件时通过邮件、钉钉等通知你
//          +------------------------+
//         |    应用程序/服务        |
//         |   (暴露 /metrics 接口)  |
//         +------------------------+
//                     ↑
//                     |
//             Prometheus 定时抓取
//                     |
//         +------------------------+
//         |    Prometheus Server   |
//         |   - 拉取 Metrics        |
//         |   - 存储时间序列数据     |
//         |   - PromQL 查询         |
//         +------------------------+
//                     |
//         +-----------+------------+
//         |                        |
//+------------------+   +-------------------+
//|  Alertmanager    |   |  Grafana           |
//|  (触发告警)       |   |  (绘图、仪表盘展示) |
//+------------------+   +-------------------+

// Readonly volume types
const (
	Namespace        = "Weedfilesys"
	IsReadOnly       = "IsReadOnly"
	NoWriteOrDelete  = "noWriteOrDelete"
	NoWriteCanDelete = "noWriteCanDelete"
	IsDiskSpaceLow   = "isDiskSpaceLow"
	bucketAtiveTTL   = 10 * time.Minute
)

var readOnlyVolumeTypes = [4]string{IsReadOnly, NoWriteOrDelete, NoWriteCanDelete, IsDiskSpaceLow}

var bucketLastActiveTsNs map[string]int64 = map[string]int64{}
var bucketLastActiveLock sync.Mutex

var (
	Gather = prometheus.NewRegistry()

	// MasterClientConnectCounter 在 Weedfilesys 中注册一个 Prometheus 指标：MasterClientConnectCounter，它用于统计 master 客户端连接（或者 leader 更新）的次数。
	// 这是指标的元信息，具体如下：
	// 字段	 			含义
	// Namespace	指标前缀，本项目是 "Weedfilesys"，用于归类所属系统
	// Subsystem	子系统名称，这里是 "wdclient"，代表 weed client
	// Name			指标名 "connect_updates"，完整名字会变成 Weedfilesys_wdclient_connect_updates
	// Help			指标说明，Prometheus UI 会显示它，用于理解这个指标干什么用的

	// 这个计数器还支持一个 标签维度：type，你可以用来区分不同类型的连接事件。例如：
	// SeaweedFS_wdclient_connect_updates{type="leader_change"}  5
	// SeaweedFS_wdclient_connect_updates{type="heartbeat"}      12
	MasterClientConnectCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "wdclient",
			Name:      "connect_updates",
			Help:      "Counter of master client leader updates.",
		}, []string{"type"})
)

func init() {
	Gather.MustRegister(MasterClientConnectCounter)
}

// 将采集到的指标数据（Gather 注册器里的指标）周期性地推送到 Prometheus PushGateway。
// | 参数                | 含义                                         |
// | -----------------  | ------------------------------------------ |
// | `name`             | job 名称，会传给 PushGateway 作为 job 名            |
// | `instance`         | 实例名，一般是主机名:端口，会作为分组标签                      |
// | `addr`             | PushGateway 的地址，比如 `http://localhost:9091` |
// | `intervalSeconds`  | 推送的时间间隔，单位秒                                |
func LoopPushingMetric(name, instance, addr string, intervalSeconds int) {
	if addr == "" || intervalSeconds == 0 {
		return
	}
	glog.V(0).Infof("%s server sends metrics to %s every %d seconds", name, addr, intervalSeconds)
	// 构建一个 Prometheus 的 Push 任务对象，具体含义：
	// 函数链				说明
	// push.New(addr, name)	创建一个新的 PushGateway 客户端，job 名为 name
	// .Gatherer(Gather)	指定要推送的指标注册器（默认 Gather 收集器）
	// .Grouping("instance", instance)	添加标签 {instance="xx.xx.xx.xx:port"}，用于区分多台机器或服务实例
	pusher := push.New(addr, name).Gatherer(Gather).Grouping("instance", instance)
	for {
		err := pusher.Push()
		if err != nil && !strings.HasPrefix(err.Error(), "unexpected status code 200") {
			glog.V(0).Infof("could not push metrics to prometheus push gateway %s: %v", addr, err)
		}
		if intervalSeconds <= 0 {
			intervalSeconds = 15
		}
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
	}
}
func JoinHostPort(host string, port int) string {
	portStr := strconv.Itoa(port)
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return host + ":" + portStr
	}
	return net.JoinHostPort(host, portStr)
}

// StartMetricsServer 会在指定 IP 和端口启动一个 HTTP 服务，暴露 /metrics 接口，供 Prometheus 访问采集监控指标数据。
func StartMetricsServer(ip string, port int) {
	if port == 0 {
		return
	}
	// 注册 HTTP 路由 /metrics，用 Prometheus 官方提供的 promhttp.HandlerFor 来生成 handler：
	//
	// 部分	    含义
	// /metrics	被 Prometheus 拉取指标的标准路径
	// Gather	前面定义并注册的 Prometheus 指标集合
	// HandlerOpts{}	默认配置（比如禁用 compression、日志、错误处理等）
	http.Handle("/metrics", promhttp.HandlerFor(Gather, promhttp.HandlerOpts{}))
	glog.Fatal(http.ListenAndServe(JoinHostPort(ip, port), nil))
}
func SourceName(port uint32) string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return net.JoinHostPort(hostname, strconv.Itoa(int(port)))
}
