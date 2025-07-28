package types

import "cayoyibackend/weedfilesys/util"

// Size 表示 Needle 的数据大小（不是偏移，是实际的 Data 部分大小）。
// 定义为 int32 是为了支持 负数 作为特殊标记（如删除标记）。
type Size int32

// | 常量名                    | 说明                                   |
// | ----------------------- | ------------------------------------ |
// | `SizeSize = 4`          | size 字段大小是 4 字节（uint32）              |
// | `NeedleHeaderSize`      | 一个 needle 的头部大小 = cookie + id + size |
// | `DataSizeSize = 4`      | 表示 Data 字段长度的数据字段，也是 4 字节            |
// | `NeedleMapEntrySize`    | 在内存中索引结构的条目大小，包含 ID、偏移、大小            |
// | `TimestampSize = 8`     | 时间戳是 `int64`，8 字节                    |
// | `NeedlePaddingSize = 8` | 所有数据都对齐到 8 字节                        |
// | `TombstoneFileSize`     | 表示已删除 needle 的特殊 size 值              |
// | `CookieSize = 4`        | cookie 是 4 字节（uint32）                |
const (
	SizeSize           = 4 // uint32 size
	NeedleHeaderSize   = CookieSize + NeedleIdSize + SizeSize
	DataSizeSize       = 4
	NeedleMapEntrySize = NeedleIdSize + OffsetSize + SizeSize
	TimestampSize      = 8 // int64 size
	NeedlePaddingSize  = 8
	TombstoneFileSize  = Size(-1)
	CookieSize         = 4
)

func (s Size) IsDeleted() bool {
	return s < 0 || s == TombstoneFileSize
}

func (s Size) IsVaild() bool {
	return s > 0 && s != TombstoneFileSize
}

func BytesToSize(bytes []byte) Size {
	return Size(util.BytesToUint32(bytes))
}

func SizeToBytes(bytes []byte, size Size) {
	util.Uint32toBytes(bytes, uint32(size))
}
