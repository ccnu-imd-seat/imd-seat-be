package logic

import (
	"context"
	"errors"
	"time"

	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/contextx"
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
	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}

	//检验信誉分
	ok, err = CheckScore(l.ctx, l.svcCtx, studentID)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	} else if !ok {
		return nil, errorx.WrapError(errorx.ViolateErr, err)
	}

	format := "2006-01-02"
	t, err := time.Parse(format, req.Date)

	ReservationInfro := model.Reservation{
		StudentId: studentID,
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

// 判断信誉分是否足够
func CheckScore(ctx context.Context, svcCtx *svc.ServiceContext, StudentID string) (bool, error) {
	//查询信誉分
	score, err := svcCtx.UserModel.FindScoreByID(ctx, StudentID)
	if err != nil {
		return false, err
	}
	if score > 0 {
		return true, nil
	}
	return false, nil
}
