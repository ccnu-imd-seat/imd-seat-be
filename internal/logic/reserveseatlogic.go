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
	debug := l.ctx.Value("DEBUG_MODE")
	fmt.Println(debug)	
	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}
	format := "2006-01-02"
	t, err := time.ParseInLocation(format, req.Date, time.Local)
	if err != nil {
		return nil, errorx.WrapError(errorx.DefaultErr, errors.New("解析时间失败"))
	}

	// 检验预约时间是否符合规则
	if debug != "1" {
		if !CheckRule(t, req.Type) {
			return nil, errorx.WrapError(errorx.ViolateErr, errors.New("预约时间不符合规则"))
		}
	}

	// 检验座位是否可预约
	seat, err := l.svcCtx.SeatModel.FindOneBySeatRoomDate(l.ctx, req.SeatID, req.Room, t)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	} else if seat == nil {
		return nil, errorx.WrapError(errorx.ViolateErr, errors.New("座位不存在"))
	} else if seat.Status != types.AvaliableStatus {
		return nil, errorx.WrapError(errorx.ViolateErr, errors.New("座位已被预约"))
	}

	// 检验信誉分
	ok, err = CheckScore(l.ctx, l.svcCtx, studentID)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	} else if !ok {
		return nil, errorx.WrapError(errorx.ViolateErr, err)
	}

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

	// 更新座位状态为已预约
	err = l.svcCtx.SeatModel.ChangeSeatStatusByType(l.ctx, t, types.BookedStatus, req.SeatID, req.Type)
	if err != nil {
		return nil, errorx.WrapError(errorx.UpdateErr, err)
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

// 检验预约时间是否符合规则
func CheckRule(date time.Time, types string) bool {
	now := time.Now()
	weekday := int(now.Weekday())
	if types == "day" {
		daysUntilSunday := (7 - weekday) % 7
		sunday := now.AddDate(0, 0, daysUntilSunday)
		if date.Before(sunday) && InTimeRange(now, 18, 21) {
			return true
		}
	} else if types == "week" {
		if weekday == 0 && InTimeRange(now, 9, 18) {
			return true
		}
	}
	return false
}

func InTimeRange(t time.Time, startHour, EndHour int) bool {
	hour := t.Hour()
	return hour >= startHour && hour <= EndHour
}
