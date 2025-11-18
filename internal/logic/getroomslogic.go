package logic

import (
	"context"

	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
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
	Rooms, err := l.svcCtx.RoomModel.GetAvailableRoom(l.ctx, types.AvaliableStatus)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}
	var roominfro []string
	for _, room := range Rooms {
		roominfro = append(roominfro, room.Room)
	}
	resp = &types.RoomListRes{
		Base: response.Success(),
		Data: types.RoomList{
			Rooms: roominfro,
		},
	}
	return
}
