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

type CheckInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInLogic {
	return &CheckInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckInLogic) CheckIn(req *types.CheckIn) (resp *types.GeneralRes, err error) {
	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}

	reservation, err := l.svcCtx.ReservationModel.GetTodayReservationByStudentId(l.ctx, studentID, req.Seatid)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}

	if reservation.Status == types.EffectiveStatus {
		return nil, errorx.WrapError(errorx.AlreadyErr, errors.New("请勿重复签到"))
	} else if reservation.Status == types.BookedStatus {
		reservation.Status = types.EffectiveStatus
		err := l.svcCtx.ReservationModel.Update(l.ctx, reservation)
		if err != nil {
			return nil, errorx.WrapError(errorx.UpdateErr, err)
		}
	} else {
		return nil, errorx.WrapError(errorx.NonCheckErr, errors.New("请预约座位"))
	}

	resp = response.GeneralRes()

	return resp, nil
}
