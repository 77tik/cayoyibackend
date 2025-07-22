package job

import (
	"fmt"
	"net/http"
	"time"

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
			filename := fmt.Sprintf("jobs-%s.zip", time.Now().Format("20060102-150405"))

			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}
	}
}
