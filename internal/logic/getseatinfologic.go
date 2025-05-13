package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSeatInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSeatInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSeatInfoLogic {
	return &GetSeatInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSeatInfoLogic) GetSeatInfo() (resp *types.SeatListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
