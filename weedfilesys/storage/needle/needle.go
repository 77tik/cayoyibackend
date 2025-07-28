package needle

import (
	"bytes"
	"cayoyibackend/weedfilesys/images"
	. "cayoyibackend/weedfilesys/storage/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	NeedleChecksumSize = 4              //每个needle都有一个4字节的CRC32校验值
	PairNamePrefix     = "Weedfilesys-" // 上传时用于过滤元数据字段的前缀，内部只保留去掉该前缀后的字段
)

// 在weedfilesys中一个上传的文件会被打包为Needle，这是在磁盘上的基本单位
type Needle struct {
	Cookie Cookie   `comment:"random number to ignore force lookup"`
	Id     NeedleId `comment:"needle id"`
	Size   Size     `comment:"sum of DataSize,Data,NameSize,Name,MimeSize,Mime"` // 后续字段的长度

	DataSize     uint32 `comment:"the real Data s Size"`
	Data         []byte `comment:"the real Data bytes"`
	Flags        byte   `comment:"binary flag"`
	NameSize     uint8  //version2
	Name         []byte `comment:"maximum 255 characters"` //version2
	MimeSize     uint8  //version2
	Mime         []byte `comment:"maximum 255 characters"` //version2
	PairsSize    uint16 //version2
	Pairs        []byte `comment:"additional name value pairs, json format, maximum 64kB"`
	LastModified uint64 //only store LastModifiedBytesLength bytes, which is 5 bytes to disk
	Ttl          *TTL

	Checksum   CRC    `comment:"CRC32 to check integrity"`
	AppendAtNs uint64 `comment:"append timestamp in nano seconds"` //version3
	Padding    []byte `comment:"Aligned to 8 bytes"`
}

// 返回 needle 的字符串描述
func (n *Needle) String() (str string) {
	str = fmt.Sprintf("%s Size:%d, DataSize:%d, Name:%s, Mime:%s Compressed:%v", formatNeedleIdCookie(n.Id, n.Cookie), n.Size, n.DataSize, n.Name, n.Mime, n.IsCompressed())
	return
}

// 从 HTTP 请求中构建一个 Needle
func CreateNeedleFromRequest(r *http.Request, fixJpgOrientation bool, sizeLimit int64, bytesBuffer *bytes.Buffer) (n *Needle, originalSize int, contentMd5 string, e error) {
	n = new(Needle)

	// 解析 HTTP 请求上传内容
	pu, e := ParseUpload(r, sizeLimit, bytesBuffer)
	if e != nil {
		return
	}

	// 设置文件数据
	n.Data = pu.Data
	originalSize = pu.OriginalDataSize
	n.LastModified = pu.ModifiedTime
	n.Ttl = pu.Ttl
	contentMd5 = pu.ContentMd5

	// 设置文件名标志位
	if len(pu.FileName) < 256 {
		n.Name = []byte(pu.FileName)
		n.SetHasName()
	}

	// 设置 Mime 类型标志位
	if len(pu.MimeType) < 256 {
		n.Mime = []byte(pu.MimeType)
		n.SetHasMime()
	}

	// 设置扩展字段 Pairs 标志位
	if len(pu.PairMap) != 0 {
		trimmedPairMap := make(map[string]string)
		for k, v := range pu.PairMap {
			trimmedPairMap[k[len(PairNamePrefix):]] = v
		}
		pairs, _ := json.Marshal(trimmedPairMap)
		if len(pairs) < 65536 {
			n.Pairs = pairs
			n.PairsSize = uint16(len(pairs))
			n.SetHasPairs()
		}
	}

	// 是否是 gzip 压缩文件，设置flag压缩标志位
	if pu.IsGzipped {
		n.SetIsCompressed()
	}

	// 如果没有设置时间戳，则使用当前时间
	if n.LastModified == 0 {
		n.LastModified = uint64(time.Now().Unix())
	}
	n.SetHasLastModifiedDate()

	// 是否设置 TTL标志位
	if n.Ttl != EMPTY_TTL {
		n.SetHasTtl()
	}

	// 如果是 分块的就设置分块标志位
	if pu.IsChunkedFile {
		n.SetIsChunkManifest()
	}

	// 是否修复 JPEG 图片方向
	if fixJpgOrientation {
		loweredName := strings.ToLower(pu.FileName)
		if pu.MimeType == "image/jpeg" || strings.HasSuffix(loweredName, ".jpg") || strings.HasSuffix(loweredName, ".jpeg") {
			n.Data = images.FixJpgOrientation(n.Data)
		}
	}

	// 生成 CRC 校验码
	n.Checksum = NewCRC(n.Data)

	// 解析文件 ID
	commaSep := strings.LastIndex(r.URL.Path, ",")
	dotSep := strings.LastIndex(r.URL.Path, ".")
	fid := r.URL.Path[commaSep+1:]
	if dotSep > 0 {
		fid = r.URL.Path[commaSep+1 : dotSep]
	}

	// 提取id和cookie填充进n
	e = n.ParsePath(fid)

	return
}

// 从 fid 中提取 NeedleId 和 Cookie 填充进n
// <needleId><cookie>_<delta>
func (n *Needle) ParsePath(fid string) (err error) {
	length := len(fid)
	if length <= CookieSize*2 {
		return fmt.Errorf("Invalid fid: %s", fid)
	}
	delta := ""
	deltaIndex := strings.LastIndex(fid, "_")
	if deltaIndex > 0 {
		fid, delta = fid[0:deltaIndex], fid[deltaIndex+1:]
	}

	n.Id, n.Cookie, err = ParseNeedleIdCookie(fid)
	if err != nil {
		return err
	}

	if delta != "" {
		if d, e := strconv.ParseUint(delta, 10, 64); e == nil {
			n.Id += Uint64ToNeedleId(d)
		} else {
			return e
		}
	}
	return err
}

// 获取当前时间与上次 Append 时间的最大值（防止覆盖）
func GetAppendAtNs(volumeLastAppendAtNs uint64) uint64 {
	return max(uint64(time.Now().UnixNano()), volumeLastAppendAtNs+1)
}

// 更新 Needle 的 AppendAtNs 时间
func (n *Needle) UpdateAppendAtNs(volumeLastAppendAtNs uint64) {
	n.AppendAtNs = max(uint64(time.Now().UnixNano()), volumeLastAppendAtNs+1)
}

// 从字符串中解析 NeedleId 和 Cookie
// key_hash_string 是把NeedleId 和Cookie 拼在一起后用16进制表示的字符串
// needleId:  0x1122334455667788  // 8字节 -> "1122334455667788" => 变成字符串就是16字节了，一个字符对应一个byte
// cookie:    0x99AABBCC           // 4字节 -> "99aabbcc"
// key_hash_string: "112233445566778899aabbcc"
func ParseNeedleIdCookie(key_hash_string string) (NeedleId, Cookie, error) {
	if len(key_hash_string) <= CookieSize*2 {
		return NeedleIdEmpty, 0, fmt.Errorf("KeyHash is too short.")
	}
	if len(key_hash_string) > (NeedleIdSize+CookieSize)*2 {
		return NeedleIdEmpty, 0, fmt.Errorf("KeyHash is too long.")
	}
	split := len(key_hash_string) - CookieSize*2
	needleId, err := ParseNeedleId(key_hash_string[:split])
	if err != nil {
		return NeedleIdEmpty, 0, fmt.Errorf("Parse needleId error: %w", err)
	}
	cookie, err := ParseCookie(key_hash_string[split:])
	if err != nil {
		return NeedleIdEmpty, 0, fmt.Errorf("Parse cookie error: %w", err)
	}
	return needleId, cookie, nil
}

// 返回 LastModified 的可读格式
func (n *Needle) LastModifiedString() string {
	return time.Unix(int64(n.LastModified), 0).Format("2006-01-02T15:04:05")
}

// 返回 x 和 y 中较大的值
func max(x, y uint64) uint64 {
	if x <= y {
		return y
	}
	return x
}
