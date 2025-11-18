package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"imd-seat-be/internal/model"
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
	var reservation *model.Reservation
	debug := l.ctx.Value("DEBUG_MODE")
	debug_day := l.ctx.Value("DEBUG_DAY")

	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}

	if debug != "1" {
		now := time.Now()
		// 获取当天 12:00 时间
		noon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
		if now.After(noon) {
			return nil, errorx.WrapError(errorx.AfterCheckErr, errors.New("无法签到"))
		}
	}

	_, err = time.Parse(time.DateOnly, debug_day.(string))

	// 这里的错误交给最下方else逻辑处理
	if debug == "1" && err == nil {
		reservation, err = l.svcCtx.ReservationModel.GetAnydayReservationByStudentId(l.ctx, studentID, req.Seatid, debug_day.(string))
	} else {
		reservation, err = l.svcCtx.ReservationModel.GetTodayReservationByStudentId(l.ctx, studentID, req.Seatid)
	}

	// 已经签到的预约单
	if reservation != nil && reservation.Status == types.EffectiveStatus {
		return nil, errorx.WrapError(errorx.AlreadyErr, errors.New("请勿重复签到"))
		// 已经预约的预约单
	} else if reservation != nil && reservation.Status == types.BookedStatus {
		reservation.Status = types.EffectiveStatus
		// 更新预约单
		err := l.svcCtx.ReservationModel.Update(l.ctx, reservation)
		if err != nil {
			return nil, errorx.WrapError(errorx.UpdateErr, err)
		}

		now := time.Now().Format("2006-01-02")
		t, err := time.ParseInLocation("2006-01-02", now, time.Local)
		if err != nil {
			return nil, errorx.WrapError(errorx.DefaultErr, errors.New("解析时间错误"))
		}

		// 更新当天座位
		err = l.svcCtx.SeatModel.ChangeSeatStatusByType(l.ctx, t, types.EffectiveStatus, reservation.Seat, "day")
		if err != nil {
			return nil, errorx.WrapError(errorx.UpdateErr, err)
		}
	} else {
		return nil, errorx.WrapError(errorx.NonCheckErr, fmt.Errorf("未预约该座位或获取数据库出错：%s", err))
	}

	resp = response.GeneralRes()

	return resp, nil
}
