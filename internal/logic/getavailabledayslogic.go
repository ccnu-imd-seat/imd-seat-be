package logic

import (
	"context"
	"fmt"
	"time"

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
	now := time.Now().Local()
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
		monday := t.AddDate(0, 0, -int(weekday)+1)
		sunday := monday.AddDate(0, 0, 6)
		datestr := fmt.Sprintf("%s——%s", monday.Format("2006-01-02"),sunday.Format("2006-01-02"))
		date.Date = datestr
		dates = append(dates, date)
	}
	return types.AvailableDates{
		Dates: dates,
	}
}
