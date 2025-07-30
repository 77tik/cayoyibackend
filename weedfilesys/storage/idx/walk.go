package idx

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/storage/types"
	"io"
)

// .idx 文件结构：
// | 字段名    | 大小（字节） | 说明                          |
// | ------ | ------ | --------------------------- |
// | key    | 8      | 文件的唯一 ID（NeedleId）          |
// | offset | 4      | 数据在 `.dat` 文件中的偏移（单位为 8 字节） |
// | size   | 4      | 实际存储的数据长度（包括 Needle 包装）     |
//
// idx 包的主要作用是：
// 对 SeaweedFS 中 .idx 索引文件进行解析、遍历、查找等操作的工具包。
// .idx 是每个 volume 文件的索引文件，记录了每个存储对象（needle）的：
// key（NeedleId）
// offset（在 .dat 文件中的位置）
// size（数据长度，包含头部和尾部）

// walks through the index file, calls fn function with each key, offset, size
// stops with the error returned by the fn function
// 这个函数的作用是遍历整个 .idx 文件内容，按 16 字节一条记录读取，每次最多读取 1024 条
func WalkIndexFile(r io.ReaderAt, startFrom uint64, fn func(key types.NeedleId, offset types.Offset, size types.Size) error) error {
	readerOffset := int64(startFrom * types.NeedleMapEntrySize)
	bytes := make([]byte, types.NeedleMapEntrySize*RowsToRead)
	count, e := r.ReadAt(bytes, readerOffset)
	if count == 0 && e == io.EOF {
		return nil
	}
	glog.V(3).Infof("readerOffset %d count %d err: %v", readerOffset, count, e)
	readerOffset += int64(count)
	var (
		key    types.NeedleId
		offset types.Offset
		size   types.Size
		i      int
	)

	for count > 0 && e == nil || e == io.EOF {
		for i = 0; i+types.NeedleMapEntrySize <= count; i += types.NeedleMapEntrySize {
			key, offset, size = IdxFileEntry(bytes[i : i+types.NeedleMapEntrySize])
			if e = fn(key, offset, size); e != nil {
				return e
			}
		}
		if e == io.EOF {
			return nil
		}
		count, e = r.ReadAt(bytes, readerOffset)
		glog.V(3).Infof("readerOffset %d count %d err: %v", readerOffset, count, e)
		readerOffset += int64(count)
	}
	return e
}

// 这是一个解码函数，用于从 16 字节的 byte 切片中解析出一条 .idx 记录的三要素：
func IdxFileEntry(bytes []byte) (key types.NeedleId, offset types.Offset, size types.Size) {
	key = types.BytesToNeedleId(bytes[:types.NeedleIdSize])
	offset = types.BytesToOffset(bytes[types.NeedleIdSize : types.NeedleIdSize+types.OffsetSize])
	size = types.BytesToSize(bytes[types.NeedleIdSize+types.OffsetSize : types.NeedleIdSize+types.OffsetSize+types.SizeSize])
	return
}

const (
	RowsToRead = 1024
)
