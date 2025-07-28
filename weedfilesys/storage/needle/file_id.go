package needle

import (
	. "cayoyibackend/weedfilesys/storage/types"
	"encoding/hex"
	"fmt"
	"strings"
)

// 暂时没看出来这个有什么用，但拥有卷和needle，看起来并不简单

// FileId 结构体定义了文件的唯一标识符
type FileId struct {
	VolumeId VolumeId // 卷ID，标识文件所在的存储卷
	Key      NeedleId // 文件键，标识文件在卷中的位置
	Cookie   Cookie   // 安全cookie，用于验证文件ID的有效性
}

// NewFileIdFromNeedle 从Needle对象创建FileId
func NewFileIdFromNeedle(VolumeId VolumeId, n *Needle) *FileId {
	return &FileId{VolumeId: VolumeId, Key: n.Id, Cookie: n.Cookie}
}

// NewFileId 使用给定的参数创建FileId
func NewFileId(VolumeId VolumeId, key uint64, cookie uint32) *FileId {
	return &FileId{
		VolumeId: VolumeId,
		Key:      Uint64ToNeedleId(key),
		Cookie:   Uint32ToCookie(cookie),
	}
}

// ParseFileIdFromString 从字符串解析FileId
func ParseFileIdFromString(fid string) (*FileId, error) {
	// 分割字符串获取卷ID和文件键+cookie部分
	vid, needleKeyCookie, err := splitVolumeId(fid)
	if err != nil {
		return nil, err
	}

	// 解析卷ID
	volumeId, err := NewVolumeId(vid)
	if err != nil {
		return nil, err
	}

	// 解析文件键和cookie
	nid, cookie, err := ParseNeedleIdCookie(needleKeyCookie)
	if err != nil {
		return nil, err
	}

	// 创建并返回FileId对象
	return &FileId{
		VolumeId: volumeId,
		Key:      nid,
		Cookie:   cookie,
	}, nil
}

// GetVolumeId 获取卷ID
func (n *FileId) GetVolumeId() VolumeId {
	return n.VolumeId
}

// GetNeedleId 获取文件键
func (n *FileId) GetNeedleId() NeedleId {
	return n.Key
}

// GetCookie 获取安全cookie
func (n *FileId) GetCookie() Cookie {
	return n.Cookie
}

// GetNeedleIdCookie 获取文件键和cookie的组合字符串表示
func (n *FileId) GetNeedleIdCookie() string {
	return formatNeedleIdCookie(n.Key, n.Cookie)
}

// String 返回FileId的字符串表示形式
func (n *FileId) String() string {
	return n.VolumeId.String() + "," + formatNeedleIdCookie(n.Key, n.Cookie)
}

// formatNeedleIdCookie 格式化文件键和cookie为十六进制字符串
func formatNeedleIdCookie(key NeedleId, cookie Cookie) string {
	// 创建足够大的字节切片来存储文件键和cookie
	bytes := make([]byte, NeedleIdSize+CookieSize)

	// 将文件键和cookie转换为字节
	NeedleIdToBytes(bytes[0:NeedleIdSize], key)
	CookieToBytes(bytes[NeedleIdSize:NeedleIdSize+CookieSize], cookie)

	// 跳过前导零字节
	nonzero_index := 0
	for ; bytes[nonzero_index] == 0 && nonzero_index < NeedleIdSize; nonzero_index++ {
	}

	// 将非零部分编码为十六进制字符串
	return hex.EncodeToString(bytes[nonzero_index:])
}

// splitVolumeId 分割文件ID字符串为卷ID和文件键+cookie部分
func splitVolumeId(fid string) (vid string, key_cookie string, err error) {
	// 查找逗号分隔符
	commaIndex := strings.Index(fid, ",")
	if commaIndex <= 0 {
		return "", "", fmt.Errorf("wrong fid format")
	}

	// 返回分割后的两部分
	return fid[:commaIndex], fid[commaIndex+1:], nil
}
