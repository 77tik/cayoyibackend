package chunk_cache

import (
	"cayoyibackend/weedfilesys/storage/backend"
	"time"
)

// ChunkCacheVolume 是 OnDiskCacheLayer 中具体负责单个磁盘文件缓存的核心实现，
// 负责管理一个磁盘缓存文件的数据读写和索引。
type ChunkCacheVolume struct {
	DataBackend backend.BackendStorageFile //封装的底层磁盘文件读写接口
	//nm          storage.NeedleMapper        // 存储 chunk 索引（needle id -> 偏移+大小）
	fileName    string    // 缓存文件名基础部分
	smallBuffer []byte    // 用于补齐对齐的临时缓冲
	sizeLimit   int64     // 文件最大大小限制
	lastModTime time.Time // 最近修改时间（用于排序管理）
	fileSize    int64     // 当前已写入的文件大小
}
