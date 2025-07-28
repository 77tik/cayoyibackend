package user

import (
	"net/http"

	"cayoyibackend/internal/logic/user"
	"cayoyibackend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取 RSA 加密公钥
func GetPubKeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetPubKeyLogic(r.Context(), svcCtx)
		resp, err := l.GetPubKey()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
