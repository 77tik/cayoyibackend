package chunk_cache

import (
	"github.com/karlseguin/ccache/v2"
	"time"
)

var (
	_ ChunkCache = &ChunkCacheInMemory{}
)

// ChunkCacheInMemory 是基于 ccache 实现的一个最近使用（LRU）策略的内存缓存层，
// 用于保存 chunk 内容并支持偏移量读取。
// 作为在内存中的缓存层
type ChunkCacheInMemory struct {
	cache *ccache.Cache //ccache.Cache 是一个并发安全、自动过期、LRU 策略的缓存结构
	// 每个缓存项是一个 fileId -> []byte 的映射，表示该 chunk 的完整数据。
}

func (c *ChunkCacheInMemory) ReadChunkAt(data []byte, fileId string, offset uint64) (n int, err error) {
	return c.readChunkAt(data, fileId, offset)
}

// 是否命中缓存，检查是否有该filed的chunk
// lockNeeded 参数被忽略，因 ccache 是线程安全的。
func (c *ChunkCacheInMemory) IsInCache(fileId string, lockNeeded bool) (answer bool) {
	item := c.cache.Get(fileId)
	if item == nil {
		return false
	}
	return true
}

func (c *ChunkCacheInMemory) GetMaxFilePartSizeInCache() (answer uint64) {
	return 8 * 1024 * 1024
}

func NewChunkCacheInMemory(maxEntries int64) *ChunkCacheInMemory {
	pruneCount := maxEntries >> 3
	if pruneCount <= 0 {
		pruneCount = 500
	}
	return &ChunkCacheInMemory{
		cache: ccache.New(ccache.Configure().MaxSize(maxEntries).ItemsToPrune(uint32(pruneCount))),
	}
}

func (c *ChunkCacheInMemory) GetChunk(fileId string) []byte {
	item := c.cache.Get(fileId)
	if item == nil {
		return nil
	}
	data := item.Value().([]byte)
	item.Extend(time.Hour)
	return data
}

// 读取整个chunk的一部分内容：
// 返回某个 chunk 的子切片；
//
// 超出边界返回 ErrorOutOfBounds；
//
// 每次访问都会 Extend 生命周期（默认续 1 小时）
func (c *ChunkCacheInMemory) getChunkSlice(fileId string, offset, length uint64) ([]byte, error) {
	item := c.cache.Get(fileId)
	if item == nil {
		return nil, nil
	}
	data := item.Value().([]byte)
	item.Extend(time.Hour)
	wanted := min(int(length), len(data)-int(offset))
	if wanted < 0 {
		return nil, ErrorOutOfBounds
	}
	return data[offset : int(offset)+wanted], nil
}

// 读取并复制到目标缓冲区（对外公开）
// 根据 offset 从缓存读取数据填充到 buffer；
//
// 会尽量拷贝完整数据，如果超过实际缓存大小就返回可读范围；
//
// 同样续命缓存。
func (c *ChunkCacheInMemory) readChunkAt(buffer []byte, fileId string, offset uint64) (int, error) {
	item := c.cache.Get(fileId)
	if item == nil {
		return 0, nil
	}
	data := item.Value().([]byte)
	item.Extend(time.Hour)
	wanted := min(len(buffer), len(data)-int(offset))
	if wanted < 0 {
		return 0, ErrorOutOfBounds
	}
	n := copy(buffer, data[offset:int(offset)+wanted])
	return n, nil
}

// 数据写入前进行一次 deep copy（拷贝到新 slice）；
// 防止外部对原始数据进行修改。
func (c *ChunkCacheInMemory) SetChunk(fileId string, data []byte) {
	localCopy := make([]byte, len(data))
	copy(localCopy, data)
	c.cache.Set(fileId, localCopy, time.Hour)
}
