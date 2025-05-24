package logic

import (
	"context"
	"errors"
	"time"

	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
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
	format := "2006-01-02"
	t, err := time.Parse(format, date)
	if err != nil {
		return nil, errorx.WrapError(errorx.DefaultErr, errors.New("解析时间失败"))
	}
	seatinfro, err := l.svcCtx.SeatModel.GetSeatInfobyDateAndID(l.ctx, t, room)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}
	var SeatInfro []types.SeatInfo
	for _, seat := range seatinfro {
		SeatInfro = append(SeatInfro, types.SeatInfo{
			SeatID: int(seat.Id),
			Status: seat.Status,
		})
	}
	resp = &types.SeatListRes{
		Base: response.Success(),
		Data: types.SeatListData{
			Room:  room,
			Date:  date,
			Seats: SeatInfro,
		},
	}
	return resp, nil
}
