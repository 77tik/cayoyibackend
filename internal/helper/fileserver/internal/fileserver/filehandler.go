package fileserver

import (
	"compress/gzip"
	"net/http"
	"path"
	"strings"
	"sync"
)

// 重写了go-zero 的fileserver部分，使其可以支持gzip压缩传输

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	w.Header().Set("Content-Encoding", "gzip")

	gzipWriter, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
	return &gzipResponseWriter{ResponseWriter: w, gzipWriter: gzipWriter}
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) { return grw.gzipWriter.Write(b) }

func (grw *gzipResponseWriter) Close() {
	_ = grw.gzipWriter.Close()
}

func Middleware(upath string, fs http.FileSystem) func(http.HandlerFunc) http.HandlerFunc {
	fileServer := http.FileServer(fs)
	pathWithoutTrailSlash := ensureNoTrailingSlash(upath)
	canServe := createServerChecker(upath, fs)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if canServe(r) {
				r.URL.Path = r.URL.Path[len(pathWithoutTrailSlash):]
				gzipRW := newGzipResponseWriter(w)
				fileServer.ServeHTTP(gzipRW, r)
				gzipRW.Close()
			} else {
				next(w, r)
			}
		}
	}
}

// ✅ 返回一个“带缓存”的 func(path string) bool 函数，用于高效判断一个文件是否存在于 http.FileSystem 中。
func createFileChecker(fs http.FileSystem) func(string) bool {
	var lock sync.RWMutex
	fileChecker := make(map[string]bool)
	return func(upath string) bool {
		// path.Clean() 会：
		//移除多余的 . 和 .。
		//合并多余的斜杠
		//保留路径结构，但输出统一规范
		upath = path.Clean("/" + upath)[1:]
		if len(upath) == 0 {
			// if the path is empty, we use "." to open the current directory
			upath = "."
		}

		lock.RLock()
		exist, ok := fileChecker[upath]
		lock.RUnlock()
		if ok {
			return exist
		}

		lock.Lock()
		defer lock.Unlock()

		file, err := fs.Open(upath)
		exist = err == nil
		fileChecker[upath] = exist
		if err != nil {
			return false
		}

		_ = file.Close()
		return true
	}
}

// ✅ 生成一个函数，用来判断某个请求 URL 是否应该由某个静态目录提供服务
// 假设你希望对 /static/ 下的资源启用静态文件服务。
// fs := http.Dir("./static")
// checker := createServeChecker("/static", fs)
func createServerChecker(upath string, fs http.FileSystem) func(r *http.Request) bool {
	// /static → /static/
	//
	// /assets/ → /assets/
	pathWithTrailSlash := ensureTrailingSlash(upath)

	// 这个 fileChecker(path) 会判断：
	//给定的 path（相对于静态目录 fs）是否存在文件？
	//例如：
	//fs = http.Dir("./static")
	//fileChecker("img/logo.png") => true/false
	fileChecker := createFileChecker(fs)

	// 假设：
	// upath = "/static"
	// fs = http.Dir("./static")  // 目录结构如下：
	//
	// ./static/
	// ├── index.html
	// └── css/style.css
	// 创建 checker：
	//checker := createServeChecker("/static", http.Dir("./static"))
	//现在有如下请求：
	///static/index.html	✅ 是	路径前缀是 /static/，存在文件 index.html
	///static/css/style.css	✅ 是	路径前缀是 /static/，存在文件 css/style.css
	///static/missing.txt	❌ 否	文件不存在
	///static2/logo.png	❌ 否	前缀不是 /static/
	///static	❌ 否	前缀不是 /static/（注意没有斜杠）
	//POST /static/index.html	❌ 否	不是 GET 请求
	return func(r *http.Request) bool {
		return r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, pathWithTrailSlash) &&
			fileChecker(r.URL.Path[len(pathWithTrailSlash):])
	}
}

func ensureTrailingSlash(upath string) string {
	if strings.HasSuffix(upath, "/") {
		return upath
	}
	return upath + "/"
}

func ensureNoTrailingSlash(upath string) string {
	if strings.HasSuffix(upath, "/") {
		return upath[:len(upath)-1]
	}
	return upath
}
