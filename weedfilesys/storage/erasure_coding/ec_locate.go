package erasure_coding

import (
	"cayoyibackend/weedfilesys/storage/types"
)

// 将偏移量和大小映射到 Erasure Coding 存储结构中的具体位置
// 即：确定某段数据在哪个 block（大/小块）、哪个 shard 文件（.ec00~.ec13）中的哪个偏移位置。以下是详细解释：
//  主要用途
//用于在 Erasure Coding（纠删码）编码的分块结构中定位实际存储的数据位置，特别是：
//
//LocateData(...)：输入数据在 .dat 文件中的偏移量和大小，输出一组 Interval，描述其在哪些块、偏移和 shard 文件中。
//
//Interval.ToShardIdAndOffset(...)：将 Interval 映射为某个 shard 文件的编号（ShardId）和该文件内的偏移量。
// 使用 Erasure Coding 将 .dat 文件划分为多个 shard 文件（如 .ec00 至 .ec13），共 14 份（10 data shards + 4 parity shards）。这些 shards 的结构是：
//
//按块（block）划分，一块数据在 10 个 data shard 上是行对齐的。
//
//块大小分为两种：大块（1GB）和小块（1MB）。
//
//每行 = 所有 10 个 data shard 各有一块。

type Interval struct {
	BlockIndex          int // the index of the block in either the large blocks or the small blocks
	InnerBlockOffset    int64
	Size                types.Size
	IsLargeBlock        bool // whether the block is a large block or a small block
	LargeBlockRowsCount int
}

func LocateData(largeBlockLength, smallBlockLength int64, shardDatSize int64, offset int64, size types.Size) (intervals []Interval) {
	blockIndex, isLargeBlock, nLargeBlockRows, innerBlockOffset := locateOffset(largeBlockLength, smallBlockLength, shardDatSize, offset)

	for size > 0 {
		interval := Interval{
			BlockIndex:          blockIndex,
			InnerBlockOffset:    innerBlockOffset,
			IsLargeBlock:        isLargeBlock,
			LargeBlockRowsCount: int(nLargeBlockRows),
		}

		blockRemaining := largeBlockLength - innerBlockOffset
		if !isLargeBlock {
			blockRemaining = smallBlockLength - innerBlockOffset
		}

		if int64(size) <= blockRemaining {
			interval.Size = size
			intervals = append(intervals, interval)
			return
		}
		interval.Size = types.Size(blockRemaining)
		intervals = append(intervals, interval)

		size -= interval.Size
		blockIndex += 1
		if isLargeBlock && blockIndex == interval.LargeBlockRowsCount*DataShardsCount {
			isLargeBlock = false
			blockIndex = 0
		}
		innerBlockOffset = 0

	}
	return
}

// DAT 文件 (逻辑视角)
// │
// ├── [LargeBlock 0] <- 每 block 分布在 ec00~ec09 的第 0 行
// ├── [LargeBlock 1]
// ...
// ├── [LargeBlock N]
// ├── [SmallBlock 0]
// ├── [SmallBlock 1]

// 想象你要存 10GB 的数据，SeaweedFS 会这样处理：
//
// 把数据划分为多个固定大小的 块（block），比如每块 1GB。
//
// 每个 块（block） 会再被分为 10 份 → 称为 data shards。
//
// 然后为每个 block，再额外计算出 4 个 校验 shards。
// 只要你还保留任意 10 个 shard（不管是数据还是校验 shard），就能 无损还原出原始数据块！
// 最终你得到一整行 14 个 shard：
//
// Block 编号（行）	ec00	ec01	ec02	...	ec09	ec10	ec11	ec12	ec13
// Block 0（第1行）	D0	D1	D2	...	D9	P0	P1	P2	P3
// Block 1（第2行）	D0	D1	D2	...	D9	P0
func locateOffset(largeBlockLength, smallBlockLength int64, shardDatSize int64, offset int64) (blockIndex int, isLargeBlock bool, nLargeBlockRows int64, innerBlockOffset int64) {
	largeRowSize := largeBlockLength * DataShardsCount
	nLargeBlockRows = (shardDatSize - 1) / largeBlockLength

	// if offset is within the large block area
	if offset < nLargeBlockRows*largeRowSize {
		isLargeBlock = true
		blockIndex, innerBlockOffset = locateOffsetWithinBlocks(largeBlockLength, offset)
		return
	}

	isLargeBlock = false
	offset -= nLargeBlockRows * largeRowSize
	blockIndex, innerBlockOffset = locateOffsetWithinBlocks(smallBlockLength, offset)
	return
}

func locateOffsetWithinBlocks(blockLength int64, offset int64) (blockIndex int, innerBlockOffset int64) {
	blockIndex = int(offset / blockLength)
	innerBlockOffset = offset % blockLength
	return
}

func (interval Interval) ToShardIdAndOffset(largeBlockSize, smallBlockSize int64) (ShardId, int64) {
	ecFileOffset := interval.InnerBlockOffset
	rowIndex := interval.BlockIndex / DataShardsCount
	if interval.IsLargeBlock {
		ecFileOffset += int64(rowIndex) * largeBlockSize
	} else {
		ecFileOffset += int64(interval.LargeBlockRowsCount)*largeBlockSize + int64(rowIndex)*smallBlockSize
	}
	ecFileIndex := interval.BlockIndex % DataShardsCount
	return ShardId(ecFileIndex), ecFileOffset
}
