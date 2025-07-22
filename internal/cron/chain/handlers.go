package chain

import (
	"context"
	"fmt"
	"my_backend/dao/model"
	"my_backend/internal/svc"
)

// 先写handler 再写branch，ctx context.Context, cctx *svc.ServiceContext, j *model.Job这些参数看着传
// 责任链内部用的ctx context.Context, cctx *svc.ServiceContext, j *model.Job 是driver自己的

func PrintHandler() Handler {
	return func(ctx context.Context, cctx *svc.ServiceContext, j *model.Job, next NextHandler) error {
		fmt.Println("PrintHandler")
		return next(ctx, cctx, j)
	}
}

func RangePrintHandler(n int) Handler {
	return func(ctx context.Context, cctx *svc.ServiceContext, j *model.Job, next NextHandler) error {
		for i := 0; i < n; i++ {
			fmt.Println("RangePrintHandler")
		}

		return next(ctx, cctx, j)
	}
}
