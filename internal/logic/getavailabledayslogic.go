package logic

import (
	"context"
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
	debug := l.ctx.Value("DEBUG_MODE")

	//获取可用天数
	dates, err := l.svcCtx.SeatModel.GetAvaliabledays(l.ctx)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}

	resp = &types.AvailableDatesRes{
		Base: response.Success(),
		Data: types.AvailableDates{
			Dates: SyncAvaliableday(Type, dates, debug),
		},
	}
	return resp, nil
}

// 整合日期
func SyncAvaliableday(Type string, dates []time.Time, debug any) []types.DateInfo {
	if len(dates) == 0 {
		return []types.DateInfo{}
	}

	var AvailableDates []types.DateInfo

	if Type == "week" && (canWeekReserve() || debug == "1") {
		// 先筛选满足“周一”且后面还有完整一周的日期
		// 先去除最后一天不是周日时，去掉最后一个周一
		lastDate := dates[len(dates)-1]
		lastWeekday := lastDate.Weekday()
		if lastWeekday != time.Sunday {
			// 去掉最后一个周一
			for i := len(dates) - 1; i >= 0; i-- {
				if dates[i].Weekday() == time.Monday {
					dates = dates[:i] // 去掉这个及后面所有元素
					break
				}
			}
		}

		for i := 0; i < len(dates); i++ {
			if dates[i].Weekday() == time.Monday && dates[i].After(time.Now()) {
				// 判断后面是否还有完整一周
				weekEnd := dates[i].AddDate(0, 0, 6) // 这一周的周日
				if !weekEnd.After(lastDate) {
					dateStr := dates[i].Format("2006-01-02")
					AvailableDates = append(AvailableDates, types.DateInfo{
						Type: "week",
						Date: dateStr,
					})
				}
				// 只留一个
				break
			}
		}

	} else if Type == "day" && (canDayReserve() || debug == "1") { // day类型，只返回本周剩余天（不含今天），若今天是周日不返回
		now := time.Now()
		todayWeekday := now.Weekday()

		if todayWeekday == time.Sunday {
			// 下周日
			nextSunday := now.AddDate(0, 0, 7)
			for _, d := range dates {
				if d.After(now) && !d.After(nextSunday) {
					AvailableDates = append(AvailableDates, types.DateInfo{
						Type: "day",
						Date: d.Format("2006-01-02"),
					})
				}
			}
		} else {
			// 返回dates中大于今天且属于本周（到周日）的日期
			// 本周周日
			sunday := now.AddDate(0, 0, 7-int(todayWeekday))
			for _, d := range dates {
				if d.After(now) && !d.After(sunday) {
					AvailableDates = append(AvailableDates, types.DateInfo{
						Type: "day",
						Date: d.Format("2006-01-02"),
					})
				}
			}
		}
	}
	return AvailableDates
}

func canWeekReserve() bool {
	now := time.Now()

	// 判断是否为周日
	if now.Weekday() != time.Sunday {
		return false
	}

	// 获取当前小时和分钟
	hour := now.Hour()
	minute := now.Minute()

	// 判断是否在 9:00 ~ 21:00 之间
	if hour < 9 || (hour == 21 && minute > 0) || hour > 21 {
		return false
	}

	return true
}

func canDayReserve() bool {
	now := time.Now()

	// 构造今天的 18:00 和 21:00
	start := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, now.Location())

	// 判断当前时间是否在区间内
	return now.After(start) && now.Before(end)
}
