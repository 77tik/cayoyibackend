package fileserver

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"my_backend/internal/helper/fileserver/internal/fileserver"
	"net/http"
	"unsafe"
)

// 重写go-zero的Server，把他的Router改成支持gzip的Router
type unsafeServer struct {
	_      unsafe.Pointer
	router httpx.Router
}

type fileServingRouter struct {
	httpx.Router
	middleware rest.Middleware
}

func WithFileServerGzip(path string, fs http.FileSystem) rest.RunOption {
	return func(server *rest.Server) {
		userver := (*unsafeServer)(unsafe.Pointer(server))
		userver.router = newFileServingRouter(userver.router, path, fs)
	}
}

func newFileServingRouter(router httpx.Router, path string, fs http.FileSystem) httpx.Router {
	return &fileServingRouter{
		Router:     router,
		middleware: fileserver.Middleware(path, fs),
	}
}

func (f *fileServingRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 如果不符合规则，就执行Router的ServeHTTP(因为next(w,r)了)，符合就执行middleware内部的逻辑然后直接返回
	f.middleware(f.Router.ServeHTTP)(w, r)
}
