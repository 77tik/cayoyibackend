package needle

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/stats"
	"cayoyibackend/weedfilesys/storage/backend"
	. "cayoyibackend/weedfilesys/storage/types"
	"cayoyibackend/weedfilesys/util"
	"errors"
	"fmt"
	"io"
)

// 定义标志位常量，Flags 字段用于表示 Needle 中可选字段的存在与否
const (
	FlagIsCompressed        = 0x01 // 是否已压缩
	FlagHasName             = 0x02 // 是否有文件名
	FlagHasMime             = 0x04 // 是否有 MIME 类型
	FlagHasLastModifiedDate = 0x08 // 是否有最后修改时间
	FlagHasTtl              = 0x10 // 是否有 TTL
	FlagHasPairs            = 0x20 // 是否有自定义键值对
	FlagIsChunkManifest     = 0x80 // 是否是 Chunk Manifest（大文件）

	LastModifiedBytesLength = 5 // 最后修改时间字段长度
	TtlBytesLength          = 2 // TTL 字段长度
)

var ErrorSizeMismatch = errors.New("size mismatch")
var ErrorSizeInvalid = errors.New("size invalid")

// Needle.DiskSize 根据版本返回写入磁盘的实际大小
func (n *Needle) DiskSize(version Version) int64 {
	return GetActualSize(n.Size, version)
}

// 从后端存储读取指定 offset 和 size 的 needle 二进制数据
func ReadNeedleBlob(r backend.BackendStorageFile, offset int64, size Size, version Version) (dataSlice []byte, err error) {
	dataSize := GetActualSize(size, version) // 包括 Header + Body + Tail
	dataSlice = make([]byte, int(dataSize))

	var n int
	n, err = r.ReadAt(dataSlice, offset) // 从指定 offset 读取到 dataSlice
	if err != nil && int64(n) == dataSize {
		err = nil // 全量读完则忽略错误（常见 EOF）
	}
	if err != nil {
		fileSize, _, _ := r.GetStat()
		glog.Errorf("%s read %d dataSize %d offset %d fileSize %d: %v", r.Name(), n, dataSize, offset, fileSize, err)
	}
	return dataSlice, err
}

// ReadBytes 从字节切片中解析 Needle 数据，且仅保证 n.Id 已设置
func (n *Needle) ReadBytes(bytes []byte, offset int64, size Size, version Version) (err error) {
	n.ParseNeedleHeader(bytes) // 先解析 Needle 头部
	if n.Size != size {        // 校验读取的尺寸是否匹配
		if OffsetSize == 4 && offset < int64(MaxPossibleVolumeSize) {
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorSizeMismatchOffsetSize).Inc() // 统计错误
			glog.Errorf("entry not found1: offset %d found id %x size %d, expected size %d", offset, n.Id, n.Size, size)
			return ErrorSizeMismatch
		}
		stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorSizeMismatch).Inc()
		return fmt.Errorf("entry not found: offset %d found id %x size %d, expected size %d", offset, n.Id, n.Size, size)
	}
	if version == Version1 {
		n.Data = bytes[NeedleHeaderSize : NeedleHeaderSize+size] // V1版本直接截取数据部分
	} else {
		err := n.readNeedleDataVersion2(bytes[NeedleHeaderSize : NeedleHeaderSize+int(size)]) // V2及以后版本解析更复杂的数据格式
		if err != nil && err != io.EOF {
			return err // 如果不是EOF错误，返回错误
		}
	}

	// 尾部CRC校验，如果是V3版本就读取创建时间
	err = n.readNeedleTail(bytes[NeedleHeaderSize+size:], version) // 解析尾部数据，如校验码等
	if err != nil {
		return err
	}
	return nil
}

// ReadData 从后端存储文件读取数据并填充 Needle 对象
// Offset 留给卷去操作吧
func (n *Needle) ReadData(r backend.BackendStorageFile, offset int64, size Size, version Version) (err error) {
	bytes, err := ReadNeedleBlob(r, offset, size, version) // 先读出字节
	if err != nil {
		return err
	}
	err = n.ReadBytes(bytes, offset, size, version)  // 再解析填充
	if err == ErrorSizeMismatch && OffsetSize == 4 { // 针对OffsetSize为4的特殊处理，尝试偏移加上最大卷大小
		offset = offset + int64(MaxPossibleVolumeSize)
		bytes, err = ReadNeedleBlob(r, offset, size, version)
		if err != nil {
			return err
		}
		err = n.ReadBytes(bytes, offset, size, version)
	}
	return err
}

// ParseNeedleHeader 解析 Needle 头部字段（Cookie、Id、Size）
func (n *Needle) ParseNeedleHeader(bytes []byte) {
	n.Cookie = BytesToCookie(bytes[0:CookieSize])                           // Cookie随机数
	n.Id = BytesToNeedleId(bytes[CookieSize : CookieSize+NeedleIdSize])     // NeedleId 唯一标识
	n.Size = BytesToSize(bytes[CookieSize+NeedleIdSize : NeedleHeaderSize]) // Size：数据大小
}

// readNeedleDataVersion2 解析V2版本数据体（包含DataSize、Data和非数据部分）
func (n *Needle) readNeedleDataVersion2(bytes []byte) (err error) {
	index, lenBytes := 0, len(bytes)
	if index < lenBytes {
		// 先读dataSize
		n.DataSize = util.BytesToUint32(bytes[index : index+4]) // 读取DataSize（4字节无符号整数）
		index = index + 4
		if int(n.DataSize)+index > lenBytes { // 越界检查
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorIndexOutOfRange).Inc()
			return fmt.Errorf("index out of range %d", 1)
		}
		n.Data = bytes[index : index+int(n.DataSize)] // 读取实际数据部分
		index = index + int(n.DataSize)
	}
	_, err = n.readNeedleDataVersion2NonData(bytes[index:]) // 读取剩余的非数据部分
	return
}

// readNeedleDataVersion2NonData 解析V2版本非数据字段（Flags、Name、Mime等）
// 返回读完之后的偏移
// 看上去和字段顺序不一致，但其实和结构体字段顺序不一样无所谓的，只要传的时候是按照这个顺序的就行
func (n *Needle) readNeedleDataVersion2NonData(bytes []byte) (index int, err error) {
	lenBytes := len(bytes)
	// 先读flag
	if index < lenBytes {
		n.Flags = bytes[index] // 读取标志位
		index = index + 1
	}
	if index < lenBytes && n.HasName() { // 有文件名
		n.NameSize = uint8(bytes[index]) // 读取名字长度
		index = index + 1
		if int(n.NameSize)+index > lenBytes { // 越界检查
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorIndexOutOfRange).Inc()
			return index, fmt.Errorf("index out of range %d", 2)
		}
		n.Name = bytes[index : index+int(n.NameSize)] // 读取名字字节
		index = index + int(n.NameSize)
	}
	if index < lenBytes && n.HasMime() { // 有Mime类型
		n.MimeSize = uint8(bytes[index]) // 读取Mime长度
		index = index + 1
		if int(n.MimeSize)+index > lenBytes { // 越界检查
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorIndexOutOfRange).Inc()
			return index, fmt.Errorf("index out of range %d", 3)
		}
		n.Mime = bytes[index : index+int(n.MimeSize)] // 读取Mime字节
		index = index + int(n.MimeSize)
	}
	if index < lenBytes && n.HasLastModifiedDate() { // 有最后修改时间
		if LastModifiedBytesLength+index > lenBytes { // 越界检查
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorIndexOutOfRange).Inc()
			return index, fmt.Errorf("index out of range %d", 4)
		}
		n.LastModified = util.BytesToUint64(bytes[index : index+LastModifiedBytesLength]) // 读取最后修改时间（5字节）
		index = index + LastModifiedBytesLength
	}
	if index < lenBytes && n.HasTtl() { // 有TTL
		if TtlBytesLength+index > lenBytes { // 越界检查
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorIndexOutOfRange).Inc()
			return index, fmt.Errorf("index out of range %d", 5)
		}
		n.Ttl = LoadTTLFromBytes(bytes[index : index+TtlBytesLength]) // 读取TTL信息
		index = index + TtlBytesLength
	}
	if index < lenBytes && n.HasPairs() { // 有附加键值对
		if 2+index > lenBytes { // 越界检查
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorIndexOutOfRange).Inc()
			return index, fmt.Errorf("index out of range %d", 6)
		}
		n.PairsSize = util.BytesToUint16(bytes[index : index+2]) // 读取Pairs大小（2字节）
		index += 2
		if int(n.PairsSize)+index > lenBytes { // 越界检查
			stats.VolumeServerHandlerCounter.WithLabelValues(stats.ErrorIndexOutOfRange).Inc()
			return index, fmt.Errorf("index out of range %d", 7)
		}
		end := index + int(n.PairsSize)
		n.Pairs = bytes[index:end] // 读取Pairs字节
		index = end
	}
	return index, nil
}

// ReadNeedleHeader 从存储文件读取 Needle 头部字节并解析
func ReadNeedleHeader(r backend.BackendStorageFile, version Version, offset int64) (n *Needle, bytes []byte, bodyLength int64, err error) {
	n = new(Needle) // 新建 Needle 对象

	bytes = make([]byte, NeedleHeaderSize) // 分配头部字节切片

	var count int
	count, err = r.ReadAt(bytes, offset)            // 从offset处读取头部
	if err == io.EOF && count == NeedleHeaderSize { // 如果读到EOF但数据完整，视为成功
		err = nil
	}
	if count <= 0 || err != nil { // 读取失败直接返回错误
		return nil, bytes, 0, err
	}

	n.ParseNeedleHeader(bytes)                     // 解析头部数据
	bodyLength = NeedleBodyLength(n.Size, version) // 计算数据体长度

	return
}

// ReadNeedleBody 从存储文件读取 Needle 数据体
func (n *Needle) ReadNeedleBody(r backend.BackendStorageFile, version Version, offset int64, bodyLength int64) (bytes []byte, err error) {

	if bodyLength <= 0 { // 数据长度小于等于0，直接返回nil
		return nil, nil
	}
	bytes = make([]byte, bodyLength)                     // 分配数据体切片
	readCount, err := r.ReadAt(bytes, offset)            // 读取数据体
	if err == io.EOF && int64(readCount) == bodyLength { // 读取到EOF但数据完整视为成功
		err = nil
	}
	if err != nil { // 读取错误时打印日志
		glog.Errorf("%s read %d bodyLength %d offset %d: %v", r.Name(), readCount, bodyLength, offset, err)
		return
	}

	err = n.ReadNeedleBodyBytes(bytes, version) // 解析数据体字节

	return
}

// ReadNeedleBodyBytes 根据版本解析数据体字节
func (n *Needle) ReadNeedleBodyBytes(needleBody []byte, version Version) (err error) {

	if len(needleBody) <= 0 { // 空数据体直接返回成功
		return nil
	}
	switch version {
	case Version1: // V1版本处理
		n.Data = needleBody[:n.Size]                         // 直接截取数据部分
		err = n.readNeedleTail(needleBody[n.Size:], version) // 读取尾部
	case Version2, Version3: // V2、V3版本处理
		err = n.readNeedleDataVersion2(needleBody[0:n.Size]) // 读取数据体内容
		if err == nil {
			err = n.readNeedleTail(needleBody[n.Size:], version) // 读取尾部
		}
	default:
		err = fmt.Errorf("unsupported version %d!", version) // 不支持的版本错误
	}
	return
}

// IsCompressed 判断数据是否压缩（根据标志位）
func (n *Needle) IsCompressed() bool {
	return n.Flags&FlagIsCompressed > 0
}

// SetIsCompressed 设置数据压缩标志
func (n *Needle) SetIsCompressed() {
	n.Flags = n.Flags | FlagIsCompressed
}

// HasName 判断是否有文件名
func (n *Needle) HasName() bool {
	return n.Flags&FlagHasName > 0
}

// SetHasName 设置文件名标志
func (n *Needle) SetHasName() {
	n.Flags = n.Flags | FlagHasName
}

// HasMime 判断是否有Mime类型
func (n *Needle) HasMime() bool {
	return n.Flags&FlagHasMime > 0
}

// SetHasMime 设置Mime标志
func (n *Needle) SetHasMime() {
	n.Flags = n.Flags | FlagHasMime
}

// HasLastModifiedDate 判断是否有最后修改时间
func (n *Needle) HasLastModifiedDate() bool {
	return n.Flags&FlagHasLastModifiedDate > 0
}

// SetHasLastModifiedDate 设置最后修改时间标志
func (n *Needle) SetHasLastModifiedDate() {
	n.Flags = n.Flags | FlagHasLastModifiedDate
}

// HasTtl 判断是否有TTL信息
func (n *Needle) HasTtl() bool {
	return n.Flags&FlagHasTtl > 0
}

// SetHasTtl 设置TTL标志
func (n *Needle) SetHasTtl() {
	n.Flags = n.Flags | FlagHasTtl
}

// IsChunkedManifest 判断是否是分块清单
func (n *Needle) IsChunkedManifest() bool {
	return n.Flags&FlagIsChunkManifest > 0
}

// SetIsChunkManifest 设置分块清单标志
func (n *Needle) SetIsChunkManifest() {
	n.Flags = n.Flags | FlagIsChunkManifest
}

// HasPairs 判断是否有附加键值对
func (n *Needle) HasPairs() bool {
	return n.Flags&FlagHasPairs != 0
}

// SetHasPairs 设置附加键值对标志
func (n *Needle) SetHasPairs() {
	n.Flags = n.Flags | FlagHasPairs
}

// GetActualSize 根据版本计算 Needle 总大小（头部 + 数据体长度）
func GetActualSize(size Size, version Version) int64 {
	return NeedleHeaderSize + NeedleBodyLength(size, version)
}
