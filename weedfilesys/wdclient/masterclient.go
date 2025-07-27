package wdclient

import (
	"context"
	"fmt"
	"my_backend/weedfilesys/util"
	"my_backend/weedfilesys/util/version"

	"google.golang.org/grpc"
	"math/rand"
	"my_backend/weedfilesys/glog"
	"my_backend/weedfilesys/pb"
	"my_backend/weedfilesys/pb/master_pb"
	"my_backend/weedfilesys/stats"
	//"my_backend/weedfilesys/util/version"
	"sync"
	"time"
)

// MasterClient 是 客户端连接 master 节点的管理器，它主要负责：
//
// 功能							实现
// 🔗 连接并保持与 master 的连接	KeepConnectedToMaster()，内含双向 gRPC 流
// 🔍 查找 fileId 所属的 volume server 地址	LookupFileIdWithFallback()
// 🧠 缓存 volumeID → server 映射	内置 vidMap
// 📡 接收 master 推送的 volume 和集群更新	处理 KeepConnected stream
// 🤝 支持 leader 节点发现与切换	FindLeaderFromOtherPeers()
// 🔄 支持 fallback 缓存链	多级 vidMap 缓存结构
type MasterClient struct {
	FilerGroup        string
	clientType        string
	clientHost        pb.ServerAddress
	rack              string
	currentMaster     pb.ServerAddress
	currentMasterLock sync.RWMutex
	masters           pb.ServerDiscovery
	grpcDialOption    grpc.DialOption

	*vidMap
	vidMapCacheSize  int
	OnPeerUpdate     func(update *master_pb.ClusterNodeUpdate, startFrom time.Time)
	OnPeerUpdateLock sync.RWMutex
}

func NewMasterClient(grpcDialOption grpc.DialOption, filerGroup string, clientType string, clientHost pb.ServerAddress, clientDataCenter string, rack string, masters pb.ServerDiscovery) *MasterClient {
	return &MasterClient{
		FilerGroup:      filerGroup,
		clientType:      clientType,
		clientHost:      clientHost,
		rack:            rack,
		masters:         masters,
		grpcDialOption:  grpcDialOption,
		vidMap:          newVidMap(clientDataCenter),
		vidMapCacheSize: 5,
	}
}
func (mc *MasterClient) SetOnPeerUpdateFn(onPeerUpdate func(update *master_pb.ClusterNodeUpdate, startFrom time.Time)) {
	mc.OnPeerUpdateLock.Lock()
	mc.OnPeerUpdate = onPeerUpdate
	mc.OnPeerUpdateLock.Unlock()
}
func (mc *MasterClient) GetLookupFileIdFunction() LookupFileIdFunctionType {
	return mc.LookupFileIdWithFallback
}

func (mc *MasterClient) LookupFileIdWithFallback(ctx context.Context, fileId string) (fullUrls []string, err error) {
	fullUrls, err = mc.vidMap.LookupFileId(ctx, fileId)
	if err == nil && len(fullUrls) > 0 {
		return
	}

	err = pb.WithMasterClient(false, mc.GetMaster(ctx), mc.grpcDialOption, false, func(client master_pb.WeedfilesysClient) error {
		resp, err := client.LookupVolume(ctx, &master_pb.LookupVolumeRequest{
			VolumeOrFileIds: []string{fileId},
		})
		if err != nil {
			return fmt.Errorf("LookupVolume %s failed: %v", fileId, err)
		}
		for vid, vidLocation := range resp.VolumeIdLocations {
			for _, vidLoc := range vidLocation.Locations {
				loc := Location{
					Url:        vidLoc.Url,
					PublicUrl:  vidLoc.PublicUrl,
					GrpcPort:   int(vidLoc.GrpcPort),
					DataCenter: vidLoc.DataCenter,
				}
				mc.vidMap.addLocation(uint32(vid), loc)
				httpUrl := "http://" + loc.Url + "/" + fileId
				// Prefer same data center
				if mc.DataCenter != "" && mc.DataCenter == loc.DataCenter {
					fullUrls = append([]string{httpUrl}, fullUrls...)
				} else {
					fullUrls = append(fullUrls, httpUrl)
				}
			}
		}
		return nil
	})
	return
}
func (mc *MasterClient) getCurrentMaster() pb.ServerAddress {
	mc.currentMasterLock.RLock()
	defer mc.currentMasterLock.RUnlock()
	return mc.currentMaster
}
func (mc *MasterClient) setCurrentMaster(master pb.ServerAddress) {
	mc.currentMasterLock.Lock()
	mc.currentMaster = master
	mc.currentMasterLock.Unlock()
}
func (mc *MasterClient) GetMaster(ctx context.Context) pb.ServerAddress {
	mc.WaitUntilConnected(ctx)
	return mc.getCurrentMaster()
}

func (mc *MasterClient) GetMasters(ctx context.Context) []pb.ServerAddress {
	mc.WaitUntilConnected(ctx)
	return mc.masters.GetInstances()
}

func (mc *MasterClient) WaitUntilConnected(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			glog.V(0).Infof("Connection wait stopped: %v", ctx.Err())
			return
		default:
			if mc.getCurrentMaster() != "" {
				return
			}
			time.Sleep(time.Duration(rand.Int31n(200)) * time.Millisecond)
			print(".")
		}
	}
}
func (mc *MasterClient) KeepConnectedToMaster(ctx context.Context) {
	glog.V(1).Infof("%s.%s masterClient bootstraps with masters %v", mc.FilerGroup, mc.clientType, mc.masters)
	for {
		select {
		case <-ctx.Done():
			glog.V(0).Infof("Connection to masters stopped: %v", ctx.Err())
			return
		default:
			mc.tryAllMasters(ctx)
			time.Sleep(time.Second)
		}
	}
}
func (mc *MasterClient) FindLeaderFromOtherPeers(myMasterAddress pb.ServerAddress) (leader string) {
	for _, master := range mc.masters.GetInstances() {
		if master == myMasterAddress {
			continue
		}
		if grpcErr := pb.WithMasterClient(false, master, mc.grpcDialOption, false, func(client master_pb.WeedfilesysClient) error {
			ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
			defer cancel()
			resp, err := client.GetMasterConfiguration(ctx, &master_pb.GetMasterConfigurationRequest{})
			if err != nil {
				return err
			}
			leader = resp.Leader
			return nil
		}); grpcErr != nil {
			glog.V(0).Infof("connect to %s: %v", master, grpcErr)
		}
		if leader != "" {
			glog.V(0).Infof("existing leader is %s", leader)
			return
		}
	}
	glog.V(0).Infof("No existing leader found!")
	return
}
func (mc *MasterClient) tryAllMasters(ctx context.Context) {
	var nextHintedLeader pb.ServerAddress
	mc.masters.RefreshBySrvIfAvailable()
	for _, master := range mc.masters.GetInstances() {
		nextHintedLeader = mc.tryConnectToMaster(ctx, master)
		for nextHintedLeader != "" {
			select {
			case <-ctx.Done():
				glog.V(0).Infof("Connection attempt to all masters stopped: %v", ctx.Err())
				return
			default:
				nextHintedLeader = mc.tryConnectToMaster(ctx, nextHintedLeader)
			}
		}
		mc.setCurrentMaster("")
	}
}

// 持续连接指定 master 节点，
// 并维护与其的 stream 通信，获取 volume 和集群更新信息，自动识别并跳转至真正的 leader 节点。
// 参数说明：
// ctx context.Context：可控制生命周期（超时/取消）；
//
// master pb.ServerAddress：尝试连接的 master 节点地址；
//
// 返回值 nextHintedLeader：如果发现当前 master 不是 leader，会返回推荐的 leader 地址，方便 caller 下次直接连对的人。
func (mc *MasterClient) tryConnectToMaster(ctx context.Context, master pb.ServerAddress) (nextHintedLeader pb.ServerAddress) {
	glog.V(1).Infof("%s.%s masterClient Connecting to master %v", mc.FilerGroup, mc.clientType, master)
	stats.MasterClientConnectCounter.WithLabelValues("total").Inc()
	grpcErr := pb.WithMasterClient(true, master, mc.grpcDialOption, false, func(client master_pb.WeedfilesysClient) (err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		stream, err := client.KeepConnected(ctx)
		if err != nil {
			glog.V(1).Infof("%s.%s masterClient failed to keep connected to %s: %v", mc.FilerGroup, mc.clientType, master, err)
			stats.MasterClientConnectCounter.WithLabelValues(stats.FailedToKeepConnected).Inc()
			return err
		}
		if err = stream.Send(&master_pb.KeepConnectedRequest{
			FilerGroup:    mc.FilerGroup,
			DataCenter:    mc.DataCenter,
			Rack:          mc.rack,
			ClientType:    mc.clientType,
			ClientAddress: string(mc.clientHost),
			Version:       version.Version(),
		}); err != nil {
			glog.V(0).Infof("%s.%s masterClient failed to send to %s: %v", mc.FilerGroup, mc.clientType, master, err)
			stats.MasterClientConnectCounter.WithLabelValues(stats.FailedToSend).Inc()
			return err
		}
		glog.V(1).Infof("%s.%s masterClient Connected to %v", mc.FilerGroup, mc.clientType, master)

		resp, err := stream.Recv()
		if err != nil {
			glog.V(0).Infof("%s.%s masterClient failed to receive from %s: %v", mc.FilerGroup, mc.clientType, master, err)
			stats.MasterClientConnectCounter.WithLabelValues(stats.FailedToReceive).Inc()
			return err
		}
		// 首次接收的时候，如果发现master不是leader，那就返回新的leader nextHintedLeader
		if resp.VolumeLocation != nil {
			if resp.VolumeLocation.Leader != "" && string(master) != resp.VolumeLocation.Leader {
				glog.V(0).Infof("master %v redirected to leader %v", master, resp.VolumeLocation.Leader)
				nextHintedLeader = pb.ServerAddress(resp.VolumeLocation.Leader)
				stats.MasterClientConnectCounter.WithLabelValues(stats.RedirectedToLeader).Inc()
				return nil
			}
			// 如果master任然是leader，那么就更新vidmap缓存
			mc.resetVidMap()
			mc.updateVidMap(resp)
		} else {
			mc.resetVidMap()
		}
		mc.setCurrentMaster(master)

		// 进入长连接持续接收master发来的更新消息
		for {
			resp, err := stream.Recv()
			if err != nil {
				glog.V(0).Infof("%s.%s masterClient failed to receive from %s: %v", mc.FilerGroup, mc.clientType, master, err)
				stats.MasterClientConnectCounter.WithLabelValues(stats.FailedToReceive).Inc()
				return err
			}

			if resp.VolumeLocation != nil {
				// maybe the leader is changed
				if resp.VolumeLocation.Leader != "" && string(mc.GetMaster(ctx)) != resp.VolumeLocation.Leader {
					glog.V(0).Infof("currentMaster %v redirected to leader %v", mc.GetMaster(ctx), resp.VolumeLocation.Leader)
					nextHintedLeader = pb.ServerAddress(resp.VolumeLocation.Leader)
					stats.MasterClientConnectCounter.WithLabelValues(stats.RedirectedToLeader).Inc()
					return nil
				}
				mc.updateVidMap(resp)
			}

			if resp.ClusterNodeUpdate != nil {
				update := resp.ClusterNodeUpdate
				mc.OnPeerUpdateLock.RLock()
				if mc.OnPeerUpdate != nil {
					if update.FilerGroup == mc.FilerGroup {
						if update.IsAdd {
							glog.V(0).Infof("+ %s@%s noticed %s.%s %s\n", mc.clientType, mc.clientHost, update.FilerGroup, update.NodeType, update.Address)
						} else {
							glog.V(0).Infof("- %s@%s noticed %s.%s %s\n", mc.clientType, mc.clientHost, update.FilerGroup, update.NodeType, update.Address)
						}
						stats.MasterClientConnectCounter.WithLabelValues(stats.OnPeerUpdate).Inc()
						mc.OnPeerUpdate(update, time.Now())
					}
				}
				mc.OnPeerUpdateLock.RUnlock()
			}
			if err := ctx.Err(); err != nil {
				glog.V(0).Infof("Connection attempt to master stopped: %v", err)
				return err
			}
		}
	})
	if grpcErr != nil {
		stats.MasterClientConnectCounter.WithLabelValues(stats.Failed).Inc()
		glog.V(1).Infof("%s.%s masterClient failed to connect with master %v: %v", mc.FilerGroup, mc.clientType, master, grpcErr)
	}
	return
}
func (mc *MasterClient) updateVidMap(resp *master_pb.KeepConnectedResponse) {
	if resp.VolumeLocation.IsEmptyUrl() {
		glog.V(0).Infof("updateVidMap ignore short heartbeat: %+v", resp)
		return
	}
	// process new volume location
	loc := Location{
		Url:        resp.VolumeLocation.Url,
		PublicUrl:  resp.VolumeLocation.PublicUrl,
		DataCenter: resp.VolumeLocation.DataCenter,
		GrpcPort:   int(resp.VolumeLocation.GrpcPort),
	}
	for _, newVid := range resp.VolumeLocation.NewVids {
		glog.V(2).Infof("%s.%s: %s masterClient adds volume %d", mc.FilerGroup, mc.clientType, loc.Url, newVid)
		mc.addLocation(newVid, loc)
	}
	for _, deletedVid := range resp.VolumeLocation.DeletedVids {
		glog.V(2).Infof("%s.%s: %s masterClient removes volume %d", mc.FilerGroup, mc.clientType, loc.Url, deletedVid)
		mc.deleteLocation(deletedVid, loc)
	}
	for _, newEcVid := range resp.VolumeLocation.NewEcVids {
		glog.V(2).Infof("%s.%s: %s masterClient adds ec volume %d", mc.FilerGroup, mc.clientType, loc.Url, newEcVid)
		mc.addEcLocation(newEcVid, loc)
	}
	for _, deletedEcVid := range resp.VolumeLocation.DeletedEcVids {
		glog.V(2).Infof("%s.%s: %s masterClient removes ec volume %d", mc.FilerGroup, mc.clientType, loc.Url, deletedEcVid)
		mc.deleteEcLocation(deletedEcVid, loc)
	}
	glog.V(1).Infof("updateVidMap(%s) %s.%s: %s volume add: %d, del: %d, add ec: %d del ec: %d",
		resp.VolumeLocation.DataCenter, mc.FilerGroup, mc.clientType, loc.Url,
		len(resp.VolumeLocation.NewVids), len(resp.VolumeLocation.DeletedVids),
		len(resp.VolumeLocation.NewEcVids), len(resp.VolumeLocation.DeletedEcVids))
}
func (mc *MasterClient) WithClient(streamingMode bool, fn func(client master_pb.WeedfilesysClient) error) error {
	getMasterF := func() pb.ServerAddress { return mc.GetMaster(context.Background()) }
	return mc.WithClientCustomGetMaster(getMasterF, streamingMode, fn)
}

func (mc *MasterClient) WithClientCustomGetMaster(getMasterF func() pb.ServerAddress, streamingMode bool, fn func(client master_pb.WeedfilesysClient) error) error {
	return util.Retry("master grpc", func() error {
		return pb.WithMasterClient(streamingMode, getMasterF(), mc.grpcDialOption, false, func(client master_pb.WeedfilesysClient) error {
			return fn(client)
		})
	})
}

func (mc *MasterClient) resetVidMap() {
	tail := &vidMap{
		vid2Locations:   mc.vid2Locations,
		ecVid2Locations: mc.ecVid2Locations,
		DataCenter:      mc.DataCenter,
		cache:           mc.cache,
	}

	nvm := newVidMap(mc.DataCenter)
	nvm.cache = tail
	mc.vidMap = nvm

	//trim
	for i := 0; i < mc.vidMapCacheSize && tail.cache != nil; i++ {
		if i == mc.vidMapCacheSize-1 {
			tail.cache = nil
		} else {
			tail = tail.cache
		}
	}
}
