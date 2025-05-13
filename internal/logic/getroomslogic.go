package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoomsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoomsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoomsLogic {
	return &GetRoomsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoomsLogic) GetRooms() (resp *types.RoomListRes, err error) {
	// todo: add your logic here and delete this line

	return
}
