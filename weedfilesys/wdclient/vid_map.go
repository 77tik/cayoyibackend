package wdclient

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"my_backend/weedfilesys/glog"
	"my_backend/weedfilesys/pb"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// 一个客户端模块，用于在本地缓存和查询 volume ID（卷 ID）对应的数据节点（location）信息，
// 从而实现快速、可用性优先的文件访问。
const (
	maxCursorIndex = 4096
)

type LookupFileIdFunctionType func(ctx context.Context, fileId string) (targetUrls []string, err error)

// vidMap 是一个“volume ID → 数据服务器地址（Location）”的本地缓存结构
type vidMap struct {
	sync.RWMutex
	vid2Locations   map[uint32][]Location
	ecVid2Locations map[uint32][]Location //支持 EC（Erasure Coding）卷的查找和缓存
	DataCenter      string                // 当前客户端所在的数据中心，用于优先选择本机房
	cursor          int32                 //原子轮询指针，用于负载平衡，避免每次都选择第一个副本
	cache           *vidMap               // 上游缓存，用于二级查询？
}

type Location struct {
	Url        string `json:"url,omitempty"`
	PublicUrl  string `json:"publicUrl,omitempty"`
	DataCenter string `json:"dataCenter,omitempty"`
	GrpcPort   int    `json:"grpcPort,omitempty"`
}

func (l Location) ServerAddress() pb.ServerAddress {
	return pb.NewServerAddressWithGrpcPort(l.Url, l.GrpcPort)
}

func newVidMap(dataCenter string) *vidMap {
	return &vidMap{
		vid2Locations:   make(map[uint32][]Location),
		ecVid2Locations: make(map[uint32][]Location),
		DataCenter:      dataCenter,
		cursor:          -1,
	}
}

func (vc *vidMap) getLocationIndex(length int) (int, error) {
	if length <= 0 {
		return 0, fmt.Errorf("invalid length: %d", length)
	}
	if atomic.LoadInt32(&vc.cursor) == maxCursorIndex {
		atomic.CompareAndSwapInt32(&vc.cursor, maxCursorIndex, -1)
	}
	return int(atomic.AddInt32(&vc.cursor, 1)) % length, nil
}
func (vc *vidMap) isSameDataCenter(loc *Location) bool {
	if vc.DataCenter == "" || loc.DataCenter == "" || vc.DataCenter != loc.DataCenter {
		return false
	}
	return true
}

func (vc *vidMap) LookupVolumeServerUrl(vid string) (serverUrls []string, err error) {
	id, err := strconv.Atoi(vid)
	if err != nil {
		glog.V(1).Infof("Unknown volume id %s", vid)
		return nil, err
	}
	locations, found := vc.GetLocations(uint32(id))
	if !found {
		return nil, fmt.Errorf("volume %d not found", id)
	}
	var sameDcServers, otherDcServers []string
	// 优先挑选与当前客户端在同一个数据中心的url
	for _, loc := range locations {
		if vc.isSameDataCenter(&loc) {
			sameDcServers = append(sameDcServers, loc.Url)
		} else {
			otherDcServers = append(otherDcServers, loc.Url)
		}
	}
	// 打乱列表顺序，防止热点过度集中
	rand.Shuffle(len(sameDcServers), func(i, j int) {
		sameDcServers[i], sameDcServers[j] = sameDcServers[j], sameDcServers[i]
	})
	rand.Shuffle(len(otherDcServers), func(i, j int) {
		otherDcServers[i], otherDcServers[j] = otherDcServers[j], otherDcServers[i]
	})

	// 优先同数据中心
	serverUrls = append(sameDcServers, otherDcServers...)
	return
}

// 解析 fileId（格式如 volumeId,hash），取出 volumeId
// 然后调用 LookupVolumeServerUrl(volumeId) 获取所有 server
// 最后拼接出完整的访问地址（如 http://ip:port/volumeId,hash）
// 🔸 用于客户端读取文件。
func (vc *vidMap) LookupFileId(ctx context.Context, fileId string) (fullUrls []string, err error) {
	parts := strings.Split(fileId, ",")
	if len(parts) != 2 {
		return nil, errors.New("Invalid fileId " + fileId)
	}
	serverUrls, lookupError := vc.LookupVolumeServerUrl(parts[0])
	if lookupError != nil {
		return nil, lookupError
	}
	for _, serverUrl := range serverUrls {
		fullUrls = append(fullUrls, "http://"+serverUrl+"/"+fileId)
	}
	return
}

func (vc *vidMap) GetVidLocations(vid string) (locations []Location, err error) {
	id, err := strconv.Atoi(vid)
	if err != nil {
		glog.V(1).Infof("Unknown volume id %s", vid)
		return nil, fmt.Errorf("Unknown volume id %s", vid)
	}
	foundLocations, found := vc.GetLocations(uint32(id))
	if found {
		return foundLocations, nil
	}
	return nil, fmt.Errorf("volume id %s not found", vid)
}
func (vc *vidMap) GetLocations(vid uint32) (locations []Location, found bool) {
	// glog.V(4).Infof("~ lookup volume id %d: %+v ec:%+v", vid, vc.vid2Locations, vc.ecVid2Locations)
	locations, found = vc.getLocations(vid)
	if found && len(locations) > 0 {
		return locations, found
	}

	if vc.cache != nil {
		return vc.cache.GetLocations(vid)
	}

	return nil, false
}

func (vc *vidMap) GetLocationsClone(vid uint32) (locations []Location, found bool) {
	locations, found = vc.GetLocations(vid)

	if found {
		// clone the locations in case the volume locations are changed below
		existingLocations := make([]Location, len(locations))
		copy(existingLocations, locations)
		return existingLocations, found
	}

	return nil, false
}
func (vc *vidMap) getLocations(vid uint32) (locations []Location, found bool) {
	vc.RLock()
	defer vc.RUnlock()

	locations, found = vc.vid2Locations[vid]
	if found && len(locations) > 0 {
		return
	}
	locations, found = vc.ecVid2Locations[vid]
	return
}

func (vc *vidMap) addLocation(vid uint32, location Location) {
	vc.Lock()
	defer vc.Unlock()

	glog.V(4).Infof("+ volume id %d: %+v", vid, location)

	locations, found := vc.vid2Locations[vid]
	if !found {
		vc.vid2Locations[vid] = []Location{location}
		return
	}

	for _, loc := range locations {
		if loc.Url == location.Url {
			return
		}
	}

	vc.vid2Locations[vid] = append(locations, location)

}

func (vc *vidMap) addEcLocation(vid uint32, location Location) {
	vc.Lock()
	defer vc.Unlock()

	glog.V(4).Infof("+ ec volume id %d: %+v", vid, location)

	locations, found := vc.ecVid2Locations[vid]
	if !found {
		vc.ecVid2Locations[vid] = []Location{location}
		return
	}

	for _, loc := range locations {
		if loc.Url == location.Url {
			return
		}
	}

	vc.ecVid2Locations[vid] = append(locations, location)

}

func (vc *vidMap) deleteLocation(vid uint32, location Location) {
	if vc.cache != nil {
		vc.cache.deleteLocation(vid, location)
	}

	vc.Lock()
	defer vc.Unlock()

	glog.V(4).Infof("- volume id %d: %+v", vid, location)

	locations, found := vc.vid2Locations[vid]
	if !found {
		return
	}

	for i, loc := range locations {
		if loc.Url == location.Url {
			vc.vid2Locations[vid] = append(locations[0:i], locations[i+1:]...)
			break
		}
	}
}

func (vc *vidMap) deleteEcLocation(vid uint32, location Location) {
	if vc.cache != nil {
		vc.cache.deleteLocation(vid, location)
	}

	vc.Lock()
	defer vc.Unlock()

	glog.V(4).Infof("- ec volume id %d: %+v", vid, location)

	locations, found := vc.ecVid2Locations[vid]
	if !found {
		return
	}

	for i, loc := range locations {
		if loc.Url == location.Url {
			vc.ecVid2Locations[vid] = append(locations[0:i], locations[i+1:]...)
			break
		}
	}
}
