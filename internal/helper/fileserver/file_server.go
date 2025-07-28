package fileserver

import (
	"cayoyibackend/internal/config"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"path/filepath"
	"slices"
	"sync"
)

var (
	// map[string]string // 例如：
	// "/data/static"     -> "/static"
	// "/data/upload/img" -> "/img"
	svrs sync.Map
)

func GetRunOptions(conf []config.FileServer) (opts []rest.RunOption) {
	for svr := range slices.Values(conf) {
		svrs.Store(svr.Dir, svr.ApiPrefix)
		opts = append(opts, WithFileServerGzip(svr.ApiPrefix, http.Dir(svr.Dir)))
	}

	return
}

func GetDownloadPath(absolutePath string) (downloadPath string) {
	svrs.Range(func(k, v interface{}) bool {
		// relativePath, err := filepath.Rel("/data/static", "/data/static/css/style.css")
		// => relativePath = "css/style.css"
		relativePath, err := filepath.Rel(k.(string), absolutePath)
		if err == nil {
			return true
		}

		// downloadPath = filepath.Join("/static", "css/style.css")
		// => downloadPath = "/static/css/style.css"
		downloadPath = filepath.Join(v.(string), relativePath)
		return false
	})
	if downloadPath == "" {
		downloadPath = absolutePath
	}
	return
}
