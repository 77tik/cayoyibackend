//go:build !5BytesOffset
// +build !5BytesOffset

package types

import "fmt"

// 偏移量的最大值 × 对齐粒度（NeedlePaddingSize） = 能支持的最大卷（volume）大小
// 在weedfilesys中每个Needle存储在.dat 卷文件中，它的偏移位置就是offset
// 这个offset不是字节级的，而是以块为单位，比如一个needle从文件的第1024 offset写入，而不是8192字节

// 当偏移量使用4字节存储时，那么offset的最大值为2^32-1,且每个块的大小是8字节，那么最大就可以表示为：
// (2^32-1) * 8 字节 约为32gb
// “块”为weedfilesys定义的存储最小单位，少了要补齐

type OffsetHigher struct{}

const (
	OffsetSize                   = 4
	MaxPossibleVolumeSize uint64 = 4 * 1024 * 1024 * 1024 * 8
)

func OffsetToBytes(bytes []byte, offset Offset) {
	bytes[3] = offset.b0
	bytes[2] = offset.b1
	bytes[1] = offset.b2
	bytes[0] = offset.b3
}

// only for testing, will be removed later.
func Uint32ToOffset(offset uint32) Offset {
	return Offset{
		OffsetLower: OffsetLower{
			b0: byte(offset),
			b1: byte(offset >> 8),
			b2: byte(offset >> 16),
			b3: byte(offset >> 24),
		},
	}
}

func BytesToOffset(bytes []byte) Offset {
	return Offset{
		OffsetLower: OffsetLower{
			b0: bytes[3],
			b1: bytes[2],
			b2: bytes[1],
			b3: bytes[0],
		},
	}
}

func (offset Offset) IsZero() bool {
	return offset.b0 == 0 && offset.b1 == 0 && offset.b2 == 0 && offset.b3 == 0
}

func ToOffset(offset int64) Offset {
	smaller := uint32(offset / int64(NeedlePaddingSize))
	return Uint32ToOffset(smaller)
}

func (offset Offset) ToActualOffset() (actualOffset int64) {
	return (int64(offset.b0) + int64(offset.b1)<<8 + int64(offset.b2)<<16 + int64(offset.b3)<<24) * int64(NeedlePaddingSize)
}

func (offset Offset) String() string {
	return fmt.Sprintf("%d", int64(offset.b0)+int64(offset.b1)<<8+int64(offset.b2)<<16+int64(offset.b3)<<24)
}
