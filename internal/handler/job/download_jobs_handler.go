package job

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"my_backend/internal/logic/job"
	"my_backend/internal/svc"
	"my_backend/internal/types"
)

// 作业文件下载
func DownloadJobsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadJobsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := job.NewDownloadJobsLogic(r.Context(), svcCtx)
		resp, err := l.DownloadJobs(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Header().Set("Content-Type", "application/zip")
			// attachment 会强制浏览器下载，filename 可根据 req 生成不同名字
			w.Header().Set("Content-Disposition", `attachment; filename="jobs.zip"`)
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}
	}
}
