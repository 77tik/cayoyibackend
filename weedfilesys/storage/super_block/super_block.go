package super_block

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/pb/master_pb"
	"cayoyibackend/weedfilesys/storage/needle"
	"cayoyibackend/weedfilesys/util"
	"google.golang.org/protobuf/proto"
)

// 🧱 什么是 SuperBlock？
// 在 SeaweedFS 的每个 Volume 数据文件开头的 前 8 字节（固定大小） 是一个 超级块（SuperBlock），用于记录与整个 Volume 相关的重要元数据。
// 它是 volume 级别的“头部”，在 volume 被创建、加载、压缩或复制时都会用到

// Byte 0   : Version （1 或 2 或 3）
// Byte 1   : Replica Placement 策略（副本数与分布策略）
// Byte 2-3 : TTL 生存时间（Time To Live）
// Byte 4-5 : Compaction Revision（卷压缩次数）
// Byte 6-7 : ExtraSize 扩展字段大小（仅 V2、V3）
// ...      : Extra 扩展字段（protobuf 序列化）
type SuperBlock struct {
	Version            needle.Version             // 卷格式版本
	ReplicaPlacement   *ReplicaPlacement          // 副本放置策略（如：001 表示一个副本在本机，一个在机架）
	Ttl                *needle.TTL                // 生存时间 TTL
	CompactionRevision uint16                     // 被压缩的次数（每次 compact 会 +1）
	Extra              *master_pb.SuperBlockExtra // 扩展字段，使用 protobuf 定义
	ExtraSize          uint16                     // 扩展字段字节长度（最大 64KB）
}

const (
	SuperBlockSize = 8
)

func (s *SuperBlock) BlockSize() int {
	switch s.Version {
	case needle.Version2, needle.Version3:
		return SuperBlockSize + int(s.ExtraSize)
	}
	return SuperBlockSize
}

func (s *SuperBlock) Bytes() []byte {
	header := make([]byte, SuperBlockSize)
	header[0] = byte(s.Version)
	header[1] = s.ReplicaPlacement.Byte()
	s.Ttl.ToBytes(header[2:4])
	util.Uint16toBytes(header[4:6], s.CompactionRevision)

	if s.Extra != nil {
		extraData, err := proto.Marshal(s.Extra)
		if err != nil {
			glog.Fatalf("cannot marshal super block extra %+v: %v", s.Extra, err)
		}
		extraSize := len(extraData)
		if extraSize > 256*256-2 {
			// reserve a couple of bits for future extension
			glog.Fatalf("super block extra size is %d bigger than %d", extraSize, 256*256-2)
		}
		s.ExtraSize = uint16(extraSize)
		util.Uint16toBytes(header[6:8], s.ExtraSize)

		header = append(header, extraData...)
	}

	return header
}

func (s *SuperBlock) Initialized() bool {
	return s.ReplicaPlacement != nil && s.Ttl != nil
}
