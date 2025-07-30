package super_block

import "fmt"

// 副本策略配置结构体 ReplicaPlacement，它决定了 每个文件副本该如何分布在不同节点、机架、数据中心中。

// ReplicaPlacement{SameRackCount: 1, DiffRackCount: 1, DiffDataCenterCount: 0}
// 表示：
// 本地保存原始副本 + 同机架上再多 1 份副本；
// 另一个机架再保存 1 份副本；
// 没有跨数据中心副本；
// 共 1（本地）+1+1 = 3 个副本。
type ReplicaPlacement struct {
	SameRackCount       int // 同一机架上额外的副本数
	DiffRackCount       int // 不同机架上的额外副本数
	DiffDataCenterCount int // 不同数据中心的额外副本数
}

// "001" => DiffDC=0, DiffRack=0, SameRack=1
// "100" => DiffDC=1, DiffRack=0, SameRack=0
// 字符串转副本结构
func NewReplicaPlacementFromString(t string) (*ReplicaPlacement, error) {
	rp := &ReplicaPlacement{}
	switch len(t) {
	case 0:
		t = "000"
	case 1:
		t = "00" + t
	case 2:
		t = "0" + t
	}
	for i, c := range t {
		count := int(c - '0')
		if count < 0 {
			return rp, fmt.Errorf("unknown replication type: %s", t)
		}
		switch i {
		case 0:
			rp.DiffDataCenterCount = count
		case 1:
			rp.DiffRackCount = count
		case 2:
			rp.SameRackCount = count
		}
	}
	value := rp.DiffDataCenterCount*100 + rp.DiffRackCount*10 + rp.SameRackCount
	if value > 255 {
		return rp, fmt.Errorf("unexpected replication type: %s", t)
	}
	return rp, nil
}

// 将单字节 byte 值（如 byte(3)）转成 ReplicaPlacement，通过先格式化为 3 位字符串再解析
func NewReplicaPlacementFromByte(b byte) (*ReplicaPlacement, error) {
	return NewReplicaPlacementFromString(fmt.Sprintf("%03d", b))
}

// 是否启用了副本机制？只要任一项 >0 就代表有副本
func (rp *ReplicaPlacement) HasReplication() bool {
	return rp.DiffDataCenterCount != 0 || rp.DiffRackCount != 0 || rp.SameRackCount != 0
}

// 比较两个副本策略是否相同。
func (a *ReplicaPlacement) Equals(b *ReplicaPlacement) bool {
	if a == nil || b == nil {
		return false
	}
	return (a.SameRackCount == b.SameRackCount &&
		a.DiffRackCount == b.DiffRackCount &&
		a.DiffDataCenterCount == b.DiffDataCenterCount)
}

func (rp *ReplicaPlacement) Byte() byte {
	if rp == nil {
		return 0
	}
	ret := rp.DiffDataCenterCount*100 + rp.DiffRackCount*10 + rp.SameRackCount
	return byte(ret)
}

func (rp *ReplicaPlacement) String() string {
	b := make([]byte, 3)
	b[0] = byte(rp.DiffDataCenterCount + '0')
	b[1] = byte(rp.DiffRackCount + '0')
	b[2] = byte(rp.SameRackCount + '0')
	return string(b)
}

// 计算实际副本总数：
// 副本总数 = 1（本地）+ SameRack + DiffRack + DiffDC
func (rp *ReplicaPlacement) GetCopyCount() int {
	return rp.DiffDataCenterCount + rp.DiffRackCount + rp.SameRackCount + 1
}
