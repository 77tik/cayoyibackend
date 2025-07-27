package filer

import "my_backend/weedfilesys/pb/filer_pb"

// Section => chunk
// 每个 Section 的大小，单位为字节。这里设为 64MiB。
const SectionSize = 2 * 1024 * 1024 * 32 // 64MiB

// SectionIndex 表示当前区块在整个文件中的第几个section
type SectionIndex int64

// FileChunkSection 表示一个大文件中的某个固定大小的区块
// 每个section区块维护自己的chunk列表，可视区间，chunk view视图和读取器
type FileChunkSection struct {
	sectionIndex     SectionIndex          //当前Section的编号
	chunks           []*filer_pb.FileChunk //属于该Section的chunk列表
	visibleIntervals *IntervalList[*VisibleInterval]
	chunkViews       *IntervalList[*ChunkView]
	reader
}
