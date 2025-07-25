package cluster

import (
	"my_backend/weedfilesys/pb"
	"time"
)

// 分组内部节点管理
type GroupMembers struct {
	members map[pb.ServerAddress]*ClusterNode
}

func newGroupMembers() *GroupMembers {
	return &GroupMembers{
		members: make(map[pb.ServerAddress]*ClusterNode),
	}
}

func (m *GroupMembers) addMember(dataCenter DataCenter, rack Rack, address pb.ServerAddress, version string) *ClusterNode {
	if existingNode, found := m.members[address]; found {
		existingNode.counter++
		return nil
	}

	t := &ClusterNode{
		Address:    address,
		DataCenter: dataCenter,
		Rack:       rack,
		Version:    version,
		counter:    1,
		CreatedTs:  time.Now(),
	}
	m.members[address] = t
	return t
}

func (m *GroupMembers) removeMember(address pb.ServerAddress) bool {
	if existingNode, found := m.members[address]; !found {
		return false
	} else {
		existingNode.counter--
		if existingNode.counter <= 0 {
			delete(m.members, address)
			return true
		}
	}
	return false
}

func (m *GroupMembers) GetMembers() (addresses []pb.ServerAddress) {
	for k := range m.members {
		addresses = append(addresses, k)
	}
	return
}
