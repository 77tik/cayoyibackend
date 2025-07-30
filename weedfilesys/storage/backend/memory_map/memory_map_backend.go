package memory_map

import (
	"os"
	"time"
)

type MemoryMappedFile struct {
	mm *MemoryMap
}

// 接收一个打开的文件 f，和最大映射大小（单位 MB）。
// 初始化一个 MemoryMap 对象，并调用 CreateMemoryMap() 将文件内容映射到内存中。
func NewMemoryMappedFile(f *os.File, memoryMapSizeMB uint32) *MemoryMappedFile {
	mmf := &MemoryMappedFile{
		mm: new(MemoryMap),
	}
	mmf.mm.CreateMemoryMap(f, 1024*1024*uint64(memoryMapSizeMB))
	return mmf
}

// 在偏移 off 处读取 len(p) 个字节。
// 调用底层 MemoryMap.ReadMemory() 从内存映射区读取数据。
// 注意：存在一次 copy()，可以优化为 zero-copy（有待 TODO）。
func (mmf *MemoryMappedFile) ReadAt(p []byte, off int64) (n int, err error) {
	readBytes, e := mmf.mm.ReadMemory(uint64(off), uint64(len(p)))
	if e != nil {
		return 0, e
	}
	copy(p, readBytes) // 把 mmap 读取的内容复制到目标缓冲区
	return len(readBytes), nil
}

// 将 p 中的字节数据写入映射文件的偏移位置 off。
// 实际是修改了内存中映射的那一段空间，写入非常快（零拷贝）。
func (mmf *MemoryMappedFile) WriteAt(p []byte, off int64) (n int, err error) {
	mmf.mm.WriteMemory(uint64(off), uint64(len(p)), p)
	return len(p), nil
}

func (mmf *MemoryMappedFile) Truncate(off int64) error {
	return nil
}
func (mmf *MemoryMappedFile) Close() error {
	mmf.mm.DeleteFileAndMemoryMap()
	return nil
}

// 获取文件状态（大小、修改时间）。
// 返回的大小为 End_of_file + 1，即当前写入结束位置。
func (mmf *MemoryMappedFile) GetStat() (datSize int64, modTime time.Time, err error) {
	stat, e := mmf.mm.File.Stat()
	if e == nil {
		return mmf.mm.End_of_file + 1, stat.ModTime(), nil
	}
	return 0, time.Time{}, err
}

func (mmf *MemoryMappedFile) Name() string {
	return mmf.mm.File.Name()
}

func (mm *MemoryMappedFile) Sync() error {
	return nil
}
