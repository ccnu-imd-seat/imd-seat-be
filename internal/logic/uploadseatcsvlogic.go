package logic

import (
	"context"

	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadSeatCsvLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadSeatCsvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadSeatCsvLogic {
	return &UploadSeatCsvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadSeatCsvLogic) UploadSeatCsv(req *types.UploadSeatRequest) (resp *types.GeneralRes, err error) {
	for _, room := range req.Rooms {
		err := l.svcCtx.SeatModel.InsertSeatsForDateRange(l.ctx, room.Room, room.Seatid, req.StartTime, req.EndTime)
		if err != nil {
			return nil, err
		}
	}

	resp = &types.GeneralRes{
		Base: response.Success(),
	}

	return resp, nil
}
