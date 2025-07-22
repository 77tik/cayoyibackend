package chain

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"my_backend/dao/model"
	"my_backend/internal/svc"
)

var JobChannel = make(chan int64, 10)

// 使用责任链启动后台任务监测Job表
func CheckJobStatus(svc *svc.ServiceContext) {
	ctx := context.Background()
	jobs := []*model.Job{&model.Job{ID: 0}}
	branchMap := map[int]*Branch{
		0: &Branch{},
	}

	go func() {
		for _, job := range jobs {
			logx.Info("CheckJobStatus jobId = %d", job.ID)
			JobChannel <- job.ID
		}
	}()

	for id := range JobChannel {
		go func() {
			b := branchMap[int(id)]
			d, err := NewDriver(
				WithSvcCtx(svc),
				WithDefaultBranch(b),
			)
			if err != nil {
				return
			}

			err = d.Chain(ctx)
			if err != nil {
				return
			}
		}()
	}
}
