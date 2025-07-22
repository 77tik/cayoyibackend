package swaggerx

import (
	"embed"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/rest"
	"io/fs"
	"net/http"
)

// 弃用了，移动到了handler的swagger handler中
var (
	//go:embed swagger-ui-5.21.0/dist
	swaggerFS embed.FS

	swaggerFSPrefix = "swagger-ui-5.21.0/dist"
)

func MustOpt() rest.RunOption {
	subfs := lo.Must(fs.Sub(swaggerFS, swaggerFSPrefix))
	return rest.WithFileServer("/swagger", http.FS(subfs))
}
