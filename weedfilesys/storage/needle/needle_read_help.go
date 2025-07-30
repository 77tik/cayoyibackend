package needle

import . "cayoyibackend/weedfilesys/storage/types"

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
