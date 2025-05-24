package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvailableDaysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvailableDaysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailableDaysLogic {
	return &GetAvailableDaysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvailableDaysLogic) GetAvailableDays(Type string) (resp *types.AvailableDatesRes, err error) {
	//获取当前时间
	now := time.Now()
	if Type == "week" {
		if now.Weekday() != time.Sunday || !InTimeRange(now, 9, 18) {
			return nil, errorx.WrapError(errorx.ViolateErr, errors.New("不在可预约时间内"))
		}
		now = now.AddDate(0, 0, 1)
	} else if Type == "day" {
		if !InTimeRange(now, 18, 21) {
			return nil, errorx.WrapError(errorx.ViolateErr, errors.New("不在可预约时间内"))
		}
	} else {
		return nil, errorx.WrapError(errorx.DefaultErr, errors.New("输入无效"))
	}
	resp = &types.AvailableDatesRes{
		Base: response.Success(),
		Data: GetThisWeekDays(now, Type),
	}
	return resp, nil
}

// 获取本周剩余日期
func GetThisWeekDays(t time.Time, Type string) types.AvailableDates {
	weekday := t.Weekday()
	if int(weekday) == 0 {
		return types.AvailableDates{
			Dates: []types.DateInfo{},
		}

	}
	DayRemaining := 7 - int(weekday)
	var date types.DateInfo
	var dates []types.DateInfo
	if Type == "day" {
		for i := 1; i <= DayRemaining; i++ {
			dates = append(dates, types.DateInfo{
				Type: "day",
				Date: t.AddDate(0, 0, i).Format("2006-01-02"),
			})
		}
	} else {
		date.Type = "week"
		datestr := fmt.Sprintf("%s——%s", t.Format("2006-01-02"), t.AddDate(0, 0, DayRemaining-1).Format("2006-01-02"))
		date.Date = datestr
		dates = append(dates, date)
	}
	return types.AvailableDates{
		Dates: dates,
	}
}

// 判断时间是否在预约范围
func InTimeRange(t time.Time, startHour, EndHour int) bool {
	hour := t.Hour()
	return hour >= startHour && hour <= EndHour
}
