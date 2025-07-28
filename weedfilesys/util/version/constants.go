package version

import (
	"cayoyibackend/weedfilesys/util"
	"fmt"
)

var (
	MAJOR_VERSION  = int32(3)                                             // 主版本号
	MINOR_VERSION  = int32(95)                                            // 次版本好
	VERSION_NUMBER = fmt.Sprintf("%d.%02d", MAJOR_VERSION, MINOR_VERSION) // 格式化为3.95 确保小数是两位
	VERSION        = util.SizeLimit + " " + VERSION_NUMBER
	COMMIT         = ""
	// COMMIT = ""
	// 用于填写构建时的 Git commit hash（构建系统可注入此值）
	//
	// 例如可以用编译参数传入：
	// go build -ldflags "-X 'version.COMMIT=abc1234'"
)

// ✅ 最终示例输出：
// 若：
//
// SizeLimit = "30GB"（从 util）
//
// MAJOR_VERSION = 3
//
// MINOR_VERSION = 95
//
// COMMIT = abc1234
//
// 调用 Version() 输出的就是：
// 30GB 3.95 abc1234
func Version() string {
	// 将版本信息组合成类似：30GB 3.95，加上体积大小限制（在 util 包中定义）
	return VERSION + " " + COMMIT
}
