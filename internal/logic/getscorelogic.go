package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScoreLogic {
	return &GetScoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScoreLogic) GetScore() (resp *types.MyScoreRes, err error) {
	// todo: add your logic here and delete this line

	return
}
