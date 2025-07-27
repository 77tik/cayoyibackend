package filer

import (
	"my_backend/weedfilesys/pb/filer_pb"
	"my_backend/weedfilesys/util"
	"os"
	"time"
)

type Attr struct {
	Mtime         time.Time   // 最后修改时间（modification time）
	Crtime        time.Time   // 创建时间（creation time，OS X 支持）
	Mode          os.FileMode // 文件模式（是否是目录、权限等）
	Uid           uint32      // 所属用户 ID
	Gid           uint32      // 所属组 ID
	Mime          string      // MIME 类型（比如 text/plain, image/jpeg）
	TtlSec        int32       // 文件的生命周期（秒）Time-To-Live
	UserName      string      // 所属用户名
	GroupNames    []string    // 所属用户组名列表
	SymlinkTarget string      // 如果是符号链接，链接的目标路径
	Md5           []byte      // 文件内容的 MD5 校验和
	FileSize      uint64      // 文件大小（单位：字节）
	Rdev          uint32      // 字符或块设备 ID（如果是设备文件）
	Inode         uint64      // inode 编号
}

func (attr Attr) IsDirectory() bool {
	return attr.Mode&os.ModeDir > 0
}

// Weedfilesys 中一个文件或者目录的表示 的基本单位
type Entry struct {
	util.FullPath // 文件路径（包含目录路径和文件名）

	Attr                       // 文件属性
	Extended map[string][]byte // 扩展属性（自定义 metadata）

	// 如果是文件，则记录其所有文件块
	Chunks []*filer_pb.FileChunk `json:"chunks,omitempty"`

	HardLinkId         HardLinkId            // 硬链接 ID（多个路径指向同一文件）
	HardLinkCounter    int32                 // 当前硬链接计数
	Content            []byte                // 直接嵌入的文件内容（小文件时使用）
	Remote             *filer_pb.RemoteEntry // 远程引用的 entry（如 cloud 存储）
	Quota              int64                 // 使用的存储配额（单位：字节）
	WORMEnforcedAtTsNs int64                 // WORM（Write Once Read Many）模式启用时间戳（纳秒）
}

// 获取当前 entry 的大小
func (entry *Entry) Size() uint64 {
	// 返回 chunks 总大小、FileSize、Content 中最大值
	return maxUint64(maxUint64(TotalSize(entry.GetChunks()), entry.FileSize), uint64(len(entry.Content)))
}

// 获取时间戳
func (entry *Entry) Timestamp() time.Time {
	if entry.IsDirectory() {
		return entry.Crtime // 目录使用创建时间
	} else {
		return entry.Mtime // 文件使用修改时间
	}
}

// 浅拷贝 只复制指针
func (entry *Entry) ShallowClone() *Entry {
	if entry == nil {
		return nil
	}
	newEntry := &Entry{}
	newEntry.FullPath = entry.FullPath
	newEntry.Attr = entry.Attr
	newEntry.Chunks = entry.Chunks
	newEntry.Extended = entry.Extended
	newEntry.HardLinkId = entry.HardLinkId
	newEntry.HardLinkCounter = entry.HardLinkCounter
	newEntry.Content = entry.Content
	newEntry.Remote = entry.Remote
	newEntry.Quota = entry.Quota

	return newEntry
}

// 构造新的 protobuf Entry 消息
func (entry *Entry) ToProtoEntry() *filer_pb.Entry {
	if entry == nil {
		return nil
	}
	message := &filer_pb.Entry{}
	message.Name = entry.FullPath.Name()
	entry.ToExistingProtoEntry(message)
	return message
}

// 用已有 protobuf Entry 结构填充内容
func (entry *Entry) ToExistingProtoEntry(message *filer_pb.Entry) {
	if entry == nil {
		return
	}
	message.IsDirectory = entry.IsDirectory()
	message.Attributes = EntryAttributeToPb(entry)
	message.Chunks = entry.GetChunks()
	message.Extended = entry.Extended
	message.HardLinkId = entry.HardLinkId
	message.HardLinkCounter = entry.HardLinkCounter
	message.Content = entry.Content
	message.RemoteEntry = entry.Remote
	message.Quota = entry.Quota
	message.WormEnforcedAtTsNs = entry.WORMEnforcedAtTsNs
}

// 从已知的protobuf entry 对象转换成本地entry对象
func FromPbEntryToExistingEntry(message *filer_pb.Entry, fsEntry *Entry) {
	fsEntry.Attr = PbToEntryAttribute(message.Attributes)
	fsEntry.Chunks = message.Chunks
	fsEntry.Extended = message.Extended
	fsEntry.HardLinkId = HardLinkId(message.HardLinkId)
	fsEntry.HardLinkCounter = message.HardLinkCounter
	fsEntry.Content = message.Content
	fsEntry.Remote = message.RemoteEntry
	fsEntry.Quota = message.Quota
	fsEntry.FileSize = FileSize(message)
	fsEntry.WORMEnforcedAtTsNs = message.WormEnforcedAtTsNs
}

// 构造fullentry，包含路径和完整内容
func (entry *Entry) ToProtoFullEntry() *filer_pb.FullEntry {
	if entry == nil {
		return nil
	}
	dir, _ := entry.FullPath.DirAndName()
	return &filer_pb.FullEntry{
		Dir:   dir,
		Entry: entry.ToProtoEntry(),
	}
}

// 获取所有块
func (entry *Entry) GetChunks() []*filer_pb.FileChunk {
	return entry.Chunks
}

// 从protobuf entry 构造一个entry，含路径
func FromPbEntry(dir string, entry *filer_pb.Entry) *Entry {
	t := &Entry{}
	t.FullPath = util.NewFullPath(dir, entry.Name)
	FromPbEntryToExistingEntry(entry, t)
	return t
}

// 返回两个uint64 中最大的那个
func maxUint64(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
