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
	resp=&types.RoomListRes{}
	Rooms,err:=l.svcCtx.RoomModel.GetAvailableRoom(l.ctx,"available")
	if err!=nil{
		resp.Base.Code=500
		resp.Base.Message="获取可预约房间失败"
		return
	}
	var roominfro []string
	for _,room:=range Rooms{
		roominfro=append(roominfro,room.Room)
	}
	resp.Base.Code=200
	resp.Base.Message="获取可预约房间成功"
	resp.Data=types.RoomList{
		Rooms: roominfro,
	}
	return
}
