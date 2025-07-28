package cluster

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/pb"
	"cayoyibackend/weedfilesys/pb/master_pb"
	"context"
	"google.golang.org/grpc"
)

//	ListExistingPeerUpdates 是 Weedfilesys 中用于查询集群中已存在的节点
//
// （如 Filer 节点）列表的逻辑，其目的是返回当前集群中某种类型（如 Filer）已注册的节点信息，
// 并包装成 ClusterNodeUpdate 结构体列表返回。
func ListExitingPeerUpdates(master pb.ServerAddress, option grpc.DialOption, filerGroup string, clientType string) (existingNodes []*master_pb.ClusterNodeUpdate) {
	if grpcErr := pb.WithMasterClient(false, master, option, false, func(client master_pb.WeedfilesysClient) error {
		resp, err := client.ListClusterNodes(context.Background(), &master_pb.ListClusterNodesRequest{
			ClientType: clientType,
			FilerGroup: filerGroup,
		})

		glog.V(0).Infof("the cluster has %d %s\n", len(resp.ClusterNodes), clientType)
		for _, node := range resp.ClusterNodes {
			existingNodes = append(existingNodes, &master_pb.ClusterNodeUpdate{
				NodeType:    FilerType,
				Address:     node.Address,
				IsAdd:       true,
				CreatedAtNs: node.CreatedAtNs,
			})
		}
		return err
	}); grpcErr != nil {
		glog.V(0).Infof("connect to %s: %v", master, grpcErr)
	}
	return
}
