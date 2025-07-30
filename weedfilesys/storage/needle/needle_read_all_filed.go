package needle

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/stats"
	"cayoyibackend/weedfilesys/storage/backend"
	. "cayoyibackend/weedfilesys/storage/types"
	"errors"
	"fmt"
	"io"
)

// 调用链关系：
// (*Needle) ReadData()
//   └── ReadNeedleBlob()  // 读原始二进制数据
//   └── (*Needle) ReadBytes()  // 解析 Needle 字段
//         └── ParseNeedleHeader()
//         └── readNeedleDataVersion2()（V2+）
//         └── readNeedleTail()

// 定义标志位常量，Flags 字段用于表示 Needle 中可选字段的存在与否
const (
	FlagIsCompressed        = 0x01 // 是否已压缩
	FlagHasName             = 0x02 // 是否有文件名
	FlagHasMime             = 0x04 // 是否有 MIME 类型
	FlagHasLastModifiedDate = 0x08 // 是否有最后修改时间
	FlagHasTtl              = 0x10 // 是否有 TTL
	FlagHasPairs            = 0x20 // 是否有自定义键值对
	FlagIsChunkManifest     = 0x80 // 是否是 Chunk Manifest（大文件）

	LastModifiedBytesLength = 5 // 最后修改时间字段长度
	TtlBytesLength          = 2 // TTL 字段长度
)

var ErrorSizeMismatch = errors.New("size mismatch")
var ErrorSizeInvalid = errors.New("size invalid")

// Needle.DiskSize 根据版本返回写入磁盘的实际大小
func (n *Needle) DiskSize(version Version) int64 {
	return GetActualSize(n.Size, version)
}

// 读二进制数据
func ReadNeedleBlob(r backend.BackendStorageFile, offset int64, size Size, version Version) (dataSlice []byte, err error) {
	dataSize := GetActualSize(size, version) // 包括 Header + Body + Tail
	dataSlice = make([]byte, int(dataSize))

	var n int
	n, err = r.ReadAt(dataSlice, offset) // 从指定 offset 读取到 dataSlice
	if err != nil && int64(n) == dataSize {
		err = nil // 全量读完则忽略错误（常见 EOF）
	}
	if err != nil {
		fileSize, _, _ := r.GetStat()
		glog.Errorf("%s read %d dataSize %d offset %d fileSize %d: %v", r.Name(), n, dataSize, offset, fileSize, err)
	}
	return dataSlice, err
}

// ReadBytes 从上面读到的二进制数据中解析 Needle 数据，且仅保证 n.Id 已设置
func (n *Needle) ReadBytes(bytes []byte, offset int64, size Size, version Version) (err error) {
	n.ParseNeedleHeader(bytes) // 先解析 Needle 头部
	if n.Size != size {        // 校验读取的尺寸是否匹配
		if OffsetSize == 4 && offset < int64(MaxPossibleVolumeSize) {
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorSizeMismatchOffsetSize).Inc() // 统计错误
			glog.Errorf("entry not found1: offset %d found id %x size %d, expected size %d", offset, n.Id, n.Size, size)
			return ErrorSizeMismatch
		}
		stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorSizeMismatch).Inc()
		return fmt.Errorf("entry not found: offset %d found id %x size %d, expected size %d", offset, n.Id, n.Size, size)
	}

	// Data这里截取了一下，还是一个二进制byte数组
	if version == Version1 {
		n.Data = bytes[NeedleHeaderSize : NeedleHeaderSize+size] // V1版本直接截取数据部分
	} else {
		err := n.readNeedleDataVersion2(bytes[NeedleHeaderSize : NeedleHeaderSize+int(size)]) // V2及以后版本解析更复杂的数据格式
		if err != nil && err != io.EOF {
			return err // 如果不是EOF错误，返回错误
		}
	}

	// 尾部CRC校验，如果是V3版本就读取创建时间
	err = n.readNeedleTail(bytes[NeedleHeaderSize+size:], version) // 解析尾部数据，如校验码等
	if err != nil {
		return err
	}
	return nil
}

// 调用上面两个读取并填充成完整的 Needle 结构体
func (n *Needle) ReadData(r backend.BackendStorageFile, offset int64, size Size, version Version) (err error) {
	bytes, err := ReadNeedleBlob(r, offset, size, version) // 先读出字节
	if err != nil {
		return err
	}
	err = n.ReadBytes(bytes, offset, size, version)  // 再解析填充
	if err == ErrorSizeMismatch && OffsetSize == 4 { // 针对OffsetSize为4的特殊处理，尝试偏移加上最大卷大小
		offset = offset + int64(MaxPossibleVolumeSize)
		bytes, err = ReadNeedleBlob(r, offset, size, version)
		if err != nil {
			return err
		}
		err = n.ReadBytes(bytes, offset, size, version)
	}
	return err
}
