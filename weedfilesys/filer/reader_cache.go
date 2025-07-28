package filer

import (
	"cayoyibackend/weedfilesys/wdclient"
	"sync"
)

// 用于缓存文件分块chunk
type ReaderCache struct {
	chunkCache     chunk_cache.ChunkCache            // 实际的 chunk 缓存实现（可用内存或磁盘缓存）
	lookupFileIdFn wdclient.LookupFileIdFunctionType // 通过 fileId 获取 volume server 地址的函数
	sync.Mutex                                       // 保护 downloaders
	downloaders    map[string]*SingleChunkCacher     // 正在下载中的 chunk
	limit          int                               // 最多同时缓存多少个 chunk
}
