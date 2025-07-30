package needle

import (
	"cayoyibackend/weedfilesys/stats"
	. "cayoyibackend/weedfilesys/storage/types"
	"cayoyibackend/weedfilesys/util"
	"fmt"
)

//// 外部解析完整 Needle 内容时：
//
//Needle.ReadData() → ReadNeedleBlob() → ReadBytes()
//                          				↓
//            				ParseNeedleHeader(bytes[:20])
//            				readNeedleDataVersion2(bytes[20:size])
//                			└── readNeedleDataVersion2NonData()

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
