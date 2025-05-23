package logic

import (
	"context"
	"errors"

	"imd-seat-be/internal/pkg/contextx"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
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
	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}

	reservations, err := l.svcCtx.ReservationModel.GetReservationByStudentId(l.ctx, studentID)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}

	var data []types.ReservationDetails

	// 结构体的转换
	for _, res := range reservations {
		data = append(data, types.ReservationDetails{
			ID:     int(res.Id),
			Type:   res.Type,
			Date:   res.Date.Format("2006-01-02"),
			Room:   res.Room,
			SeatID: res.Seat,
			Status: res.Status,
		})
	}

	resp = &types.MyReservationRes{
		Base: response.Success(),
		Data: data,
	}

	return resp, nil
}
