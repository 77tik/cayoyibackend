package filer

import "cayoyibackend/weedfilesys/wdclient"

// ChunkReadAt 负责从多个 chunk 中读取数据，是 Weedfilesys 拼接逻辑文件数据的核心组件
// 它根据 chunkViews（可见的 chunk 片段）按 offset 读取数据，可以处理稀疏区域（填充 0）
type ChunkReadAt struct {
	masterClient  *wdclient.MasterClient    // 可选，暂未用
	chunkViews    *IntervalList[*ChunkView] // 当前 section 中可读的 chunk 区间
	fileSize      int64                     // 文件总大小
	readerCache   *ReaderCache              // chunk 内容的缓存（可缓存整个 chunk 或一部分）
	readerPattern *ReaderPattern            // 记录是否是顺序读/随机读
	lastChunkFid  string                    // 用于优化缓存淘汰
}
