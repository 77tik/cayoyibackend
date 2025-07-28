package util

import (
	"bytes"
	"cayoyibackend/weedfilesys/glog"
	"fmt"
	"strings"
)

var (
	UnsupportedCompression = fmt.Errorf("unsupported compression")
)

func MaybeGzipData(input []byte) []byte {
	if IsGzippedContent(input) {
		return input // 如果已是 gzip 数据，则直接返回
	}
	gzipped, err := GzipData(input)
	if err != nil {
		return input // 压缩失败，返回原始数据
	}
	if len(gzipped)*10 > len(input)*9 {
		return input // 压缩效果差，原始数据更小，不压缩
	}
	return gzipped // 返回压缩后的数据
}

// 这是解压入口函数。会自动检测是否 gzip 格式，如果是则尝试解压，不支持的压缩格式或出错时返回原始数据
func MaybeDecompressData(input []byte) []byte {
	uncompressed, err := DecompressData(input)
	if err != nil {
		if err != UnsupportedCompression {
			glog.Errorf("decompressed data: %v", err)
		}
		return input // 解压失败就返回原始数据
	}
	return uncompressed // 返回解压后的数据
}

func GzipData(input []byte) ([]byte, error) {
	w := new(bytes.Buffer)                          // 用于存放压缩数据
	_, err := GzipStream(w, bytes.NewReader(input)) // GzipStream 是压缩流写入函数（在其他文件中实现）
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil // 返回压缩结果
}
func ungzipData(input []byte) ([]byte, error) {
	w := new(bytes.Buffer)
	_, err := GunzipStream(w, bytes.NewReader(input)) // GunzipStream 是解压流函数
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// 如果是 gzip 格式，则调用 ungzipData，否则返回一个 UnsupportedCompression 错误。
func DecompressData(input []byte) ([]byte, error) {
	if IsGzippedContent(input) {
		return ungzipData(input) // 是 gzip 格式就解压
	}
	/*
		if IsZstdContent(input) {
			return unzstdData(input)
		}
	*/
	return input, UnsupportedCompression // 不支持的格式
}

// 判断数据是否是 gzip 格式，是通过检查前两个字节是否为 gzip 的魔数（Magic Number）：1F 8B
// 在计算机文件格式中，Magic Number（魔数） 是指文件开头的特定字节序列，用于标识文件类型。
// 很多二进制格式（包括 Gzip）都会在开头写入一段固定字节作为标志。
// Gzip 格式的文件头固定前两个字节是：
// 十六进制：1F 8B
// 十进制：  31 139
// 0x1F（十进制 31）
// 0x8B（十进制 139）
func IsGzippedContent(data []byte) bool {
	if len(data) < 2 {
		return false // 数据太短，不可能是 gzip
	}
	return data[0] == 31 && data[1] == 139 // gzip magic number: 0x1f 0x8b
}

// 返回文件是否是可以压缩/已经压缩的类型
func IsCompressableFileType(ext, mtype string) (shouldBeCompressed, iAmSure bool) {
	// 文本类
	if strings.HasPrefix(mtype, "text/") {
		return true, true
	}

	// 特定可压缩的文件类型
	switch ext {
	case ".svg", ".bmp", ".wav":
		return true, true
	}

	// 图片通常已压缩，不再压缩
	if strings.HasPrefix(mtype, "image/") {
		return false, true
	}

	// 根据扩展名判断
	switch ext {
	case ".zip", ".rar", ".gz", ".bz2", ".xz", ".zst", ".br":
		return false, true // 已压缩格式
	case ".pdf", ".txt", ".html", ".htm", ".css", ".js", ".json":
		return true, true // 高压缩率文本类
	case ".php", ".java", ".go", ".rb", ".c", ".cpp", ".h", ".hpp":
		return true, true // 源代码类文件
	case ".png", ".jpg", ".jpeg":
		return false, true // 已压缩的图片格式
	}

	// 根据MIME类型进一步识别
	if strings.HasPrefix(mtype, "application/") {
		if strings.HasSuffix(mtype, "zstd") {
			return false, true
		}
		if strings.HasSuffix(mtype, "xml") {
			return true, true
		}
		if strings.HasSuffix(mtype, "script") {
			return true, true
		}
		if strings.HasSuffix(mtype, "vnd.rar") {
			return false, true
		}
	}

	// 音频压缩类型判断
	if strings.HasPrefix(mtype, "audio/") {
		switch strings.TrimPrefix(mtype, "audio/") {
		case "wave", "wav", "x-wav", "x-pn-wav":
			return true, true // wav 是未压缩音频
		}
	}

	// 兜底返回不确定
	return false, false
}
