package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"my_backend/internal/logic/user"
	"my_backend/internal/svc"
	"my_backend/internal/types"
)

// 用户信息修改
func UpdateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateUserLogic(r.Context(), svcCtx)
		err := l.UpdateUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
