package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvailableDaysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvailableDaysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailableDaysLogic {
	return &GetAvailableDaysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvailableDaysLogic) GetAvailableDays() (resp *types.AvailableDatesRes, err error) {
	// todo: add your logic here and delete this line

	return
}
