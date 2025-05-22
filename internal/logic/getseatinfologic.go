package logic

import (
	"context"
	"time"

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

func (l *GetSeatInfoLogic) GetSeatInfo(date, room string) (resp *types.SeatListRes, err error) {
	resp = &types.SeatListRes{}
	format := "2006-01-02"
	t, err := time.Parse(format, date)
	if err != nil {
		resp.Base.Code = 400
		resp.Base.Message = "日期格式错误"
		return
	}
	seatinfro, err := l.svcCtx.SeatModel.GetSeatInfobyDateAndID(l.ctx, t, room)
	if err != nil {
		resp.Base.Code = 500
		resp.Base.Message = "获取座位信息失败！"
		return
	}
	var SeatInfro []types.SeatInfo
	for _, seat := range seatinfro {
		SeatInfro = append(SeatInfro, types.SeatInfo{
			SeatID: int(seat.Id),
			Status: seat.Status,
		})
	}
	resp.Base.Code = 200
	resp.Base.Message = "获取座位数据成功!"
	resp.Data = types.SeatListData{
		Room:  room,
		Date:  date,
		Seats: SeatInfro,
	}
	return
}
