package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"my_backend/internal/logic/user"
	"my_backend/internal/svc"
)

// 用户信息
func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
