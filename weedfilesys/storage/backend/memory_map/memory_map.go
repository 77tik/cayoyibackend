package memory_map

import (
	"os"
	"strconv"
)

// 内存映射（memory mapping）文件写入的结构体定义及一个配置解析函数，主要用于对大文件进行高效读写（尤其是写入）

type MemoryBuffer struct {
	aligned_length uint64  // 分配的内存对齐长度（通常按页对齐）
	length         uint64  // 实际使用的长度（可能小于 aligned_length）
	aligned_ptr    uintptr // 对齐后的指针地址
	ptr            uintptr // 原始分配的指针地址
	Buffer         []byte  // 实际使用的字节切片
}

// 这通常用于优化写入性能，特别是当使用 mmap 技术将文件内容映射到内存中时，可以直接对文件进行内存操作而不走传统 I/O 接口。
type MemoryMap struct {
	File                   *os.File       // 映射的文件
	file_memory_map_handle uintptr        // 操作系统层的 memory map 句柄（Windows 或 syscall 层）
	write_map_views        []MemoryBuffer // 多个 memory buffer，用于映射视图（可多段）
	max_length             uint64         // 最大可映射长度
	End_of_file            int64          // 当前文件的结尾位置（可能小于 max_length）
}

// input: "512"   => 输出: 512, nil
// input: ""      => 输出: 0, nil
// input: "abc"   => 输出: 0, error
// 这个函数用于从配置参数中解析 MemoryMap 的最大使用内存限制。
func ReadMemoryMapMaxSizeMb(memoryMapMaxSizeMbString string) (uint32, error) {
	if memoryMapMaxSizeMbString == "" {
		return 0, nil
	}
	memoryMapMaxSize64, err := strconv.ParseUint(memoryMapMaxSizeMbString, 10, 32)
	return uint32(memoryMapMaxSize64), err
}
