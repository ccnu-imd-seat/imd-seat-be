package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelReservationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelReservationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelReservationLogic {
	return &CancelReservationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelReservationLogic) CancelReservation() (resp *types.GeneralRes, err error) {
	// todo: add your logic here and delete this line

	return
}
