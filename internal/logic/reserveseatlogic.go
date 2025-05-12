package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReserveSeatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReserveSeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReserveSeatLogic {
	return &ReserveSeatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReserveSeatLogic) ReserveSeat(req *types.ReserveSeatReq) (resp *types.ReserveSeatRes, err error) {
	// todo: add your logic here and delete this line

	return
}
