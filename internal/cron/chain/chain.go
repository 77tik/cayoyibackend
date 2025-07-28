package chain

import (
	"cayoyibackend/dao/model"
	"cayoyibackend/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

// 自造责任链实现定时任务
// 1个driver => 1个branch => n个handler

// 取消作业
type CancelError struct{}

func (e *CancelError) Error() string { return "cancel job" }

// 开始调用下一个任务
type NextHandler func(ctx context.Context, cctx *svc.ServiceContext, j *model.Job) error

// 用于处理请求并根据情况调用下一个处理器
type Handler func(ctx context.Context, cctx *svc.ServiceContext, j *model.Job, next NextHandler) error

// Branch 是责任链中的一个分支
type Branch struct {
	handlers []Handler
}

// 分支设置，用于配置分支参数
type BranchOptions func(*Branch)

// 配置责任链分支中的处理器
func WithBranchHandlers(handlers ...Handler) BranchOptions {
	return func(b *Branch) {
		b.handlers = append(b.handlers, handlers...)
	}
}

// 创建一个分支
func NewBranch(opts ...BranchOptions) *Branch {
	b := new(Branch)
	for _, o := range opts {
		o(b)
	}

	return b
}

type Driver struct {
	job    *model.Job
	cctx   *svc.ServiceContext
	branch *Branch
}

// 难道是为整个驱动做配置吗
type OptionOnDriver func(*Driver) error

func WithDefaultBranch(b *Branch) OptionOnDriver {
	return func(d *Driver) error {
		d.branch = b
		return nil
	}
}

func WithJob(job *model.Job) OptionOnDriver {
	return func(d *Driver) error {
		d.job = job
		return nil
	}
}

func WithSvcCtx(cctx *svc.ServiceContext) OptionOnDriver {
	return func(d *Driver) error {
		d.cctx = cctx
		return nil
	}
}

// 我还以为是一个树干开多个分叉，没想到分叉就是树干，树干就是分叉，所谓的branch看上去就是个叶子，分叉才是driver
func NewDriver(opts ...OptionOnDriver) (*Driver, error) {
	drv := &Driver{}
	for _, o := range opts {
		if err := o(drv); err != nil {
			return nil, err
		}
	}
	return drv, nil
}

func (drv *Driver) Chain(ctx context.Context) error {
	var index int
	var next NextHandler
	var currBranch = drv.branch

	//
	next = func(ctx context.Context, cctx *svc.ServiceContext, j *model.Job) error {
		if index >= len(currBranch.handlers) {
			return nil
		}
		var currHandler Handler
		currHandler, index = currBranch.handlers[index], index+1

		// 执行当前Handler，顺便传入next作为下一个Hander
		err := currHandler(ctx, cctx, j, next)
		if err != nil {
			logx.Errorf("[chain] Do Job[%d] Failed!, err = %v", j.ID, err)
			// TODO:更新数据库表状态为 取消或失败
		}

		// 别问为什么这里没看到有转移支付的代码，那个要自定义，在自定义的handler中决定是否next
		return nil
	}

	// 启动链条：
	return next(ctx, drv.cctx, drv.job)
}
