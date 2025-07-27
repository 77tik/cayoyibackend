package chunk_cache

import (
	"errors"
	"sync"
)

// 这个模块是 Weedfilesys 在读取文件时的三层级 chunk 缓存系统，
// 用于加速 chunk 的重复读取，提升性能，减少网络请求：

var ErrorOutOfBounds = errors.New("attempt to read out of bounds")

// 任何想用来做chunk的缓存的结构都必须实现整个接口才行
type ChunkCache interface {
	ReadChunkAt(data []byte, fileId string, offset uint64) (n int, err error) // 从缓存中读取 chunk 内容
	SetChunk(fileId string, data []byte)                                      // 写入 chunk 内容到缓存
	IsInCache(fileId string, lockNeeded bool) (answer bool)                   // 判断某个 chunk 是否已经缓存
	GetMaxFilePartSizeInCache() (answer uint64)                               // 当前缓存允许的最大 chunk 大小
}

type TieredChunkCache struct {
	memCache   *ChunkCacheInMemory // 内存缓存（最高速）
	diskCaches []*OnDiskCacheLayer // 三层磁盘缓存：小文件放第 0 层，大文件放第 1、2 层
	sync.RWMutex

	onDiskCacheSizeLimit0 uint64 // 层级 0 的最大 chunk 大小
	onDiskCacheSizeLimit1 uint64 // 层级 1 的最大 chunk 大小
	onDiskCacheSizeLimit2 uint64 // 层级 2 的最大 chunk 大小

	maxFilePartSizeInCache uint64 // 外部使用者使用的最大缓存大小限制
}
