package needle

import (
	"cayoyibackend/weedfilesys/glog"
	"cayoyibackend/weedfilesys/storage/backend"
	"fmt"
	"io"
)

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
