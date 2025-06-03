package logic

import (
	"context"
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
	//获取可用天数
	dates, err := l.svcCtx.SeatModel.GetAvaliabledays(l.ctx)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}
	resp = &types.AvailableDatesRes{
		Base: response.Success(),
		Data: types.AvailableDates{
			Dates: SyncAvaliableday(Type, dates),
		},
	}
	return resp, nil
}

// 整合日期
func SyncAvaliableday(Type string, dates []time.Time) []types.DateInfo {
	if len(dates)==0{
		return []types.DateInfo{}
	}
	var AvailableDates []types.DateInfo
	if Type == "week" {
		datestr1:=dates[1].Format("2006-01-02")
		datestr2:=dates[len(dates)-1].Format("2006-01-02")
		date := types.DateInfo{
			Type: "week",
			Date: fmt.Sprintf("%s - %s", datestr1,datestr2),
		}
		AvailableDates = append(AvailableDates, date)
	} else {
		for i := 1; i < len(dates); i++ {
			datestr:=dates[i].Format("2006-01-02")
			date := types.DateInfo{
				Type: "day",
				Date: datestr,
			}
			AvailableDates = append(AvailableDates, date)
		}
	}
	logx.Infof("date：%v",AvailableDates)
	return AvailableDates
}
