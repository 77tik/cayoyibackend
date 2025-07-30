package needle

import (
	"cayoyibackend/weedfilesys/storage/backend"
	. "cayoyibackend/weedfilesys/storage/types"
	"io"
)

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
