package types

import (
	"cayoyibackend/weedfilesys/util"
	"fmt"
	"strconv"
)

type Offset struct {
	OffsetHigher // 预留给5字节偏移
	OffsetLower  // 固定4字节偏移
}

// 将4字节的uint32值按小端顺序（低位在前）存储为4个字节字段，方便序列化和反序列化
type OffsetLower struct {
	b3 byte
	b2 byte
	b1 byte
	b0 byte
}

// 为什么需要 Cookie？
// 因为如果别人知道了你的 NeedleId，可能会构造访问请求去尝试获取数据。
// 而加上一个只有服务端保存的随机 Cookie，访问者必须同时提供正确的 Cookie 才能获取数据。
// 否则，即使 ID 是对的，系统也会拒绝。
type Cookie uint32

func CookieToBytes(bytes []byte, cookie Cookie) {
	util.Uint32toBytes(bytes, uint32(cookie))
}
func Uint32ToCookie(cookie uint32) Cookie {
	return Cookie(cookie)
}

func BytesToCookie(bytes []byte) Cookie {
	return Cookie(util.BytesToUint32(bytes[0:4]))
}

func ParseCookie(cookieString string) (Cookie, error) {
	cookie, err := strconv.ParseUint(cookieString, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("needle cookie %s format error: %v", cookieString, err)
	}
	return Cookie(cookie), nil
}
