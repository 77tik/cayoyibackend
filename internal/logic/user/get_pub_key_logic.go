package user

import (
	"context"

	"cayoyibackend/internal/svc"
	"cayoyibackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPubKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 RSA 加密公钥
func NewGetPubKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPubKeyLogic {
	return &GetPubKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPubKeyLogic) GetPubKey() (resp *types.GetPubKeyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
