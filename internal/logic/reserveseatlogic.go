package logic

import (
	"context"
	"errors"
	"time"

	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
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
	student_id, ok := l.ctx.Value("student_id").(string)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}
	format := "2006-01-02"
	t, err := time.Parse(format, req.Date)

	ReservationInfro := model.Reservation{
		StudentId: student_id,
		Type:      req.Type,
		Date:      t,
		Room:      req.Room,
		Seat:      req.SeatID,
		Status:    types.BookedStatus,
	}

	re, err := l.svcCtx.ReservationModel.Insert(l.ctx, &ReservationInfro)
	if err != nil {
		return nil, errorx.WrapError(errorx.CreateErr, err)
	}
	id, _ := re.LastInsertId()
	resp = &types.ReserveSeatRes{
		Base: response.Success(),
		Data: types.ReservationData{
			RoomID:        req.Room,
			Seat:          req.SeatID,
			Date:          req.Date,
			Type:          req.Type,
			ReservationID: int(id),
		},
	}
	return resp, nil
}
