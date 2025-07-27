package filer

import "my_backend/weedfilesys/pb/filer_pb"

// 计算所有 chunk 中最大的 (offset + size)，即文件的逻辑总大小
func TotalSize(chunks []*filer_pb.FileChunk) (size uint64) {
	for _, c := range chunks {
		t := uint64(c.Offset + int64(c.Size))
		if size < t {
			size = t
		}
	}
	return
}

// 根据 Entry 信息获取文件大小，考虑 chunk 和 remote entry 情况
func FileSize(entry *filer_pb.Entry) (size uint64) {
	if entry == nil || entry.Attributes == nil {
		return 0
	}
	fileSize := entry.Attributes.FileSize
	if entry.RemoteEntry != nil {
		if entry.RemoteEntry.RemoteMtime > entry.Attributes.Mtime {
			fileSize = maxUint64(fileSize, uint64(entry.RemoteEntry.RemoteSize))
		}
	}
	return maxUint64(TotalSize(entry.GetChunks()), fileSize)
}

// VisibleInterval 表示某个 chunk 在逻辑文件中的“可见范围”，用于构建有效的读取视图。
// 一个 chunk 的一部分可能被更新、覆盖或部分可见，因此必须通过该结构精确表示。
type VisibleInterval struct {
	start         int64  // 当前 chunk 可见数据在文件中的起始偏移（逻辑偏移）
	stop          int64  // 当前 chunk 可见数据在文件中的结束偏移（不包含）
	modifiedTsNs  int64  // 修改时间戳（纳秒），用于判定覆盖关系，时间新的优先
	fileId        string // 对应 chunk 的文件标识（在 Volume 中的 FileId）
	offsetInChunk int64  // 当前 chunk 的可见范围在 chunk 内部的起始偏移（即 ViewStart = chunk.Offset + offsetInChunk）
	chunkSize     uint64 // chunk 的总大小
	cipherKey     []byte // 若 chunk 加密，此字段存储加密密钥
	isGzipped     bool   // 若 chunk 被压缩，此字段为 true
}

// 调整VisibaleInterval 的逻辑可见范围，并更新offset
func (v *VisibleInterval) SetStartStop(start, stop int64) {
	// 由于 start 改变了，需要更新 chunk 内部偏移量
	v.offsetInChunk += start - v.start
	v.start, v.stop = start, stop
}

func (v *VisibleInterval) Clone() IntervalValue {
	return &VisibleInterval{
		start:         v.start,
		stop:          v.stop,
		modifiedTsNs:  v.modifiedTsNs,
		fileId:        v.fileId,
		offsetInChunk: v.offsetInChunk,
		chunkSize:     v.chunkSize,
		cipherKey:     v.cipherKey,
		isGzipped:     v.isGzipped,
	}
}

// [多个 FileChunk]
//
//	↓
//
// ResolveChunkManifest
//
//	↓
//
// [多个 VisibleInterval] ← 做版本合并、去重、处理重叠
//
//	↓
//
// ViewFromVisibleIntervals
//
//	↓
//
// [多个 ChunkView] ← 精确描述从哪个 chunk 读哪些字节
type ChunkView struct {
	FileId        string // Chunk 在 Volume 中的唯一标识符（如 "3,abc123"）
	OffsetInChunk int64  // chunk 内部偏移，从哪个位置开始读取数据
	ViewSize      uint64 // 本次读取的大小（以字节为单位）
	ViewOffset    int64  // 本段数据在逻辑文件中的偏移
	ChunkSize     uint64 // chunk 总大小（可能大于 ViewSize）
	CipherKey     []byte // 若 chunk 被加密，存储解密密钥
	IsGzipped     bool   // 是否开启 Gzip 压缩
	ModifiedTsNs  int64  // 该段数据的修改时间戳（纳秒），用于判断版本
}

func (cv *ChunkView) SetStartStop(start, stop int64) {
	cv.OffsetInChunk += start - cv.ViewOffset
	cv.ViewOffset = start
	cv.ViewSize = uint64(stop - start)
}
func (cv *ChunkView) Clone() IntervalValue {
	return &ChunkView{
		FileId:        cv.FileId,
		OffsetInChunk: cv.OffsetInChunk,
		ViewSize:      cv.ViewSize,
		ViewOffset:    cv.ViewOffset,
		ChunkSize:     cv.ChunkSize,
		CipherKey:     cv.CipherKey,
		IsGzipped:     cv.IsGzipped,
		ModifiedTsNs:  cv.ModifiedTsNs,
	}
}
func (cv *ChunkView) IsFullChunk() bool {
	return cv.ViewSize == cv.ChunkSize
}

func min(x, y int64) int64 {
	if x <= y {
		return x
	}
	return y
}

func max(x, y int64) int64 {
	if x <= y {
		return y
	}
	return x
}
