package idx

import "cayoyibackend/weedfilesys/storage/types"

// idx 包的主要作用是：
// 对 SeaweedFS 中 .idx 索引文件进行解析、遍历、查找等操作的工具包。
// .idx 是每个 volume 文件的索引文件，记录了每个存储对象（needle）的：
// key（NeedleId）
// offset（在 .dat 文件中的位置）
// size（数据长度，包含头部和尾部）

func FirstInvalidIndex(bytes []byte, lessThanOrEqualToFn func(key types.NeedleId, offset types.Offset, size types.Size) (bool, error)) (int, error) {
	left, right := 0, len(bytes)/types.NeedleMapEntrySize-1
	index := right + 1
	for left <= right {
		mid := left + (right-left)>>1
		loc := mid * types.NeedleMapEntrySize
		key := types.BytesToNeedleId(bytes[loc : loc+types.NeedleIdSize])
		offset := types.BytesToOffset(bytes[loc+types.NeedleIdSize : loc+types.NeedleIdSize+types.OffsetSize])
		size := types.BytesToSize(bytes[loc+types.NeedleIdSize+types.OffsetSize : loc+types.NeedleIdSize+types.OffsetSize+types.SizeSize])
		res, err := lessThanOrEqualToFn(key, offset, size)
		if err != nil {
			return -1, err
		}
		if res {
			left = mid + 1
		} else {
			index = mid
			right = mid - 1
		}
	}
	return index, nil
}
