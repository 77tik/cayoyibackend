package needle

import (
	"bytes"
	"cayoyibackend/weedfilesys/util"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// 从HTTP请求中解析上传的文件信息，并封装成内部结构体，供后续写入存储使用
type ParsedUpload struct {
	FileName         string            // 上传文件名（从 Content-Disposition 或 multipart 解析）
	Data             []byte            // 读取到的实际内容数据（可能被压缩）
	bytesBuffer      *bytes.Buffer     // 缓冲区，用于临时存储读取内容
	MimeType         string            // MIME 类型（从 Content-Type 或扩展名推断）
	PairMap          map[string]string // 以 Seaweed- 开头的头部键值对，用于附加信息
	IsGzipped        bool              // 是否是 gzip 压缩过的数据
	OriginalDataSize int               // 未压缩前的数据长度
	ModifiedTime     uint64            // 可选的 ts 字段，表示最后修改时间（来自表单）
	Ttl              *TTL              // 生存时间（从 URL 参数中读取，例如 ?ttl=3d）
	IsChunkedFile    bool              // 是否是 chunked 文件上传（Form 参数 cm=true）
	UncompressedData []byte            // 解压缩后的数据（等于原始数据或解压结果）
	ContentMd5       string            // base64 编码的 MD5，用于完整性校验
}

// 从HTTP请求中提取上传数据和元信息，支持表单和非表单方式上传
func ParseUpload(r *http.Request, sizeLimit int64, bytesBuffer *bytes.Buffer) (pu *ParsedUpload, e error) {
	bytesBuffer.Reset()
	pu = &ParsedUpload{bytesBuffer: bytesBuffer}
	pu.PairMap = make(map[string]string)

	// 从Header 提取以Weedfilesys- 开头的键值对
	for k, v := range r.Header {
		if len(v) > 0 && strings.HasPrefix(k, PairNamePrefix) {
			pu.PairMap[k] = v[0]
		}
	}
	// 调用辅助函数，从 HTTP 请求体中读取实际上传数据
	e = parseUpload(r, sizeLimit, pu)
	if e != nil {
		return
	}
	// 获取表单参数：ts 和 ttl
	pu.ModifiedTime, _ = strconv.ParseUint(r.FormValue("ts"), 10, 64)
	pu.Ttl, _ = ReadTTL(r.FormValue("ttl"))

	// 保存原始数据长度
	pu.OriginalDataSize = len(pu.Data)
	pu.UncompressedData = pu.Data // 默认是原始数据

	// 如果是 gzip 压缩的，尝试解压并记录原始大小
	if pu.IsGzipped {
		if unzipped, e := util.DecompressData(pu.Data); e == nil {
			pu.OriginalDataSize = len(unzipped)
			pu.UncompressedData = unzipped
		}
	} else {
		ext := filepath.Base(pu.FileName)
		mimeType := pu.MimeType
		if mimeType == "" {
			mimeType = http.DetectContentType(pu.Data) // 自动识别 MIME
		}
		if mimeType == "application/octet-stream" {
			mimeType = ""
		}
		// 判断是否是可压缩类型
		if shouldBeCompressed, iAmSure := util.IsCompressableFileType(ext, mimeType); shouldBeCompressed && iAmSure {
			if compressedData, err := util.GzipData(pu.Data); err == nil {
				if len(compressedData)*10 < len(pu.Data)*9 { // 判断压缩效果
					pu.Data = compressedData
					pu.IsGzipped = true
				}
			}
		}
	}
	// 生成 MD5 校验值
	h := md5.New()
	h.Write(pu.UncompressedData)
	pu.ContentMd5 = base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 如果请求头里有 Content-MD5 字段，就校验完整性
	if expectedChecksum := r.Header.Get("Content-MD5"); expectedChecksum != "" {
		if expectedChecksum != pu.ContentMd5 {
			e = fmt.Errorf("Content-MD5 did not match md5 of file data expected [%s] received [%s] size %d", expectedChecksum, pu.ContentMd5, len(pu.UncompressedData))
			return
		}
	}
	return
}

func parseUpload(r *http.Request, sizeLimit int64, pu *ParsedUpload) (e error) {
	// 错误时确保关闭请求体，避免连接泄漏
	defer func() {
		if e != nil && r.Body != nil {
			io.Copy(io.Discard, r.Body) // 读取剩余内容防止连接复用问题
			r.Body.Close()
		}
	}()

	contentType := r.Header.Get("Content-Type") // 获取 Content-Type
	var dataSize int64

	// 情况一：是 multipart/form-data 表单提交（比如 HTML 表单上传文件）
	if r.Method == http.MethodPost && (contentType == "" || strings.Contains(contentType, "form-data")) {
		form, fe := r.MultipartReader() // 获取 multipart 表单 reader
		if fe != nil {
			e = fe
			return
		}

		part, fe := form.NextPart() // 读取第一个 part，通常是文件字段
		if fe != nil {
			e = fe
			return
		}

		// 获取文件名
		pu.FileName = part.FileName()
		if pu.FileName != "" {
			pu.FileName = path.Base(pu.FileName) // 去掉路径，仅保留文件名
		}

		// 从 multipart 读取文件数据到内存缓冲区（带上限）
		dataSize, e = pu.bytesBuffer.ReadFrom(io.LimitReader(part, sizeLimit+1))
		if e != nil {
			return
		}
		if dataSize == sizeLimit+1 {
			e = fmt.Errorf("file over the limited %d bytes", sizeLimit)
			return
		}
		pu.Data = pu.bytesBuffer.Bytes() // 设置上传的原始数据内容

		contentType = part.Header.Get("Content-Type")

		// 如果第一个 part 没有文件名，继续读取下一个 part 尝试获取
		for pu.FileName == "" {
			part2, fe := form.NextPart()
			if fe != nil {
				break // 没有更多 part 了
			}

			fName := part2.FileName()
			if fName != "" {
				pu.bytesBuffer.Reset()
				dataSize2, fe2 := pu.bytesBuffer.ReadFrom(io.LimitReader(part2, sizeLimit+1))
				if fe2 != nil {
					e = fe2
					return
				}
				if dataSize2 == sizeLimit+1 {
					e = fmt.Errorf("file over the limited %d bytes", sizeLimit)
					return
				}
				pu.Data = pu.bytesBuffer.Bytes()
				pu.FileName = path.Base(fName)
				contentType = part.Header.Get("Content-Type")
				part = part2
				break
			}
		}

		// 判断是否设置了 gzip 编码
		pu.IsGzipped = part.Header.Get("Content-Encoding") == "gzip"

	} else {
		// 情况二：不是 multipart 表单提交，可能是 curl 上传（单纯的 body）

		// 从 Content-Disposition 中尝试解析文件名
		disposition := r.Header.Get("Content-Disposition")
		if strings.Contains(disposition, "name=") {
			if !strings.HasPrefix(disposition, "inline") && !strings.HasPrefix(disposition, "attachment") {
				disposition = "attachment; " + disposition
			}
			_, mediaTypeParams, err := mime.ParseMediaType(disposition)
			if err == nil {
				dpFilename, hasFilename := mediaTypeParams["filename"]
				dpName, hasName := mediaTypeParams["name"]
				if hasFilename {
					pu.FileName = dpFilename
				} else if hasName {
					pu.FileName = dpName
				}
			}
		} else {
			pu.FileName = ""
		}

		// 没有解析出文件名则使用 URL 路径作为文件名
		if pu.FileName != "" {
			pu.FileName = path.Base(pu.FileName)
		} else {
			pu.FileName = path.Base(r.URL.Path)
		}

		// 限制读取 body 的大小
		dataSize, e = pu.bytesBuffer.ReadFrom(io.LimitReader(r.Body, sizeLimit+1))
		if e != nil {
			return
		}
		if dataSize == sizeLimit+1 {
			e = fmt.Errorf("file over the limited %d bytes", sizeLimit)
			return
		}

		pu.Data = pu.bytesBuffer.Bytes()
		pu.MimeType = contentType
		pu.IsGzipped = r.Header.Get("Content-Encoding") == "gzip"
	}

	// 判断是否是分块上传（用于大文件切片）
	pu.IsChunkedFile, _ = strconv.ParseBool(r.FormValue("cm"))

	// 如果不是分块上传，则尝试自动推断 MIME 类型
	if !pu.IsChunkedFile {
		dotIndex := strings.LastIndex(pu.FileName, ".")
		ext, mtype := "", ""
		if dotIndex > 0 {
			ext = strings.ToLower(pu.FileName[dotIndex:])
			mtype = mime.TypeByExtension(ext)
		}

		if contentType != "" && contentType != "application/octet-stream" && mtype != contentType {
			pu.MimeType = contentType
		} else if mtype != "" && pu.MimeType == "" && mtype != "application/octet-stream" {
			pu.MimeType = mtype
		}
	}

	return
}
