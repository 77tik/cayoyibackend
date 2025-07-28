package needle

import (
	"cayoyibackend/weedfilesys/util"
	"fmt"
	"hash/crc32"
	"io"
)

// 实现 Needle 文件完整性校验的核心逻辑，它定义了 CRC 类型以及相关操作，
// 用于计算和跟踪文件数据的校验值，确保写入磁盘的数据没有损坏。
// CRC 的本质：
// 把数据当成一个二进制多项式
// 用一个“约定好”的生成多项式 G(x) 对数据做除法
// 把余数 R(x) 附加到数据尾部，发送出去
// 接收端用同样的 G(x) 再除一次：
// 如果余数为 0，说明数据没被动过
// 如果余数不为 0，说明有误

// Q: 为什么不用 MD5/SHA1 呢？
// 因为 CRC32 更快（速度几百倍），足够检测损坏。
//
// SeaweedFS 的用途是检测磁盘错误而不是防攻击，CRC 足够了。
//
// Q: 写入顺序是先数据再 CRC 吗？
// 是的，固定结构是：
// [DataSize][Data][Flags等元信息][CRC][Padding]

// 生成查找表
var table = crc32.MakeTable(crc32.Castagnoli)

type CRC uint32

func NewCRC(b []byte) CRC { return CRC(0).Update(b) }

func (c CRC) Update(b []byte) CRC { return CRC(crc32.Update(uint32(c), table, b)) }

// Value Deprecated. Just use the raw uint32 value to compare.
func (c CRC) Value() uint32 {
	return uint32(c>>15|c<<17) + 0xa282ead8
}

func (n *Needle) Etag() string {
	bits := make([]byte, 4)
	util.Uint32toBytes(bits, uint32(n.Checksum))
	return fmt.Sprintf("%x", bits)
}

func NewCRCwriter(w io.Writer) *CRCwriter {

	return &CRCwriter{
		crc: CRC(0),
		w:   w,
	}

}

type CRCwriter struct {
	crc CRC
	w   io.Writer
}

func (c *CRCwriter) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p) // with each write ...
	c.crc = c.crc.Update(p)
	return
}

func (c *CRCwriter) Sum() uint32 { return uint32(c.crc) } // final hash
