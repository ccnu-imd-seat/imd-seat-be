package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyReservationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyReservationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyReservationLogic {
	return &GetMyReservationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyReservationLogic) GetMyReservation() (resp *types.MyReservationRes, err error) {
	

	l.svcCtx.ReservationModel.GetReservationByStudentId(l.ctx,)

	return
}
