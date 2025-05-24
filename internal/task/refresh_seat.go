package task

import (
	"context"
	"log"
	"time"
	"errors"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"
	"imd-seat-be/internal/pkg/errorx"

	"github.com/robfig/cron/v3"
)

func RegisterTasks(ctx context.Context, svcCtx *svc.ServiceContext) {
	c := cron.New(cron.WithSeconds())
	//每天十二点进行更新
	_, err := c.AddFunc("0 12 0 * * *", func() {
		if err := Violation(ctx, svcCtx); err != nil {
			log.Println("更新预约状态失败:", err)
		} else {
			log.Println("更新预约状态成功")
		}
	})
	if err != nil {
		log.Println("注册定时任务失败:", err)
	}
	c.Start()

	go func() {
		<-ctx.Done()
		c.Stop()
	}()
}

// 查找所有未签到的预约,更新状态为违约
// TODO:释放座位，以及按周预约时的违约判定
func Violation(ctx context.Context, svcCtx *svc.ServiceContext) error {
	now := time.Now().Format("2006-01-02")
	parsedTime, err := time.Parse("2006-01-02", now)
	if err != nil {
		return errorx.WrapError(errorx.DefaultErr,errors.New("解析时间错误"))
	}
	Reservations, err := svcCtx.ReservationModel.GetReservationByStatus(ctx, parsedTime, types.EffectiveStatus)
	if err != nil {
		return errorx.WrapError(errorx.FetchErr,err)
	}
	var FailedIds []int64
	for _, reservation := range Reservations {
		err := svcCtx.ReservationModel.UpdateReservstionMessage(ctx, reservation.Id, types.ViolatedStatus)
		if err != nil {
			log.Printf("更新id:%d的预约状态失败%v", reservation.Id, err)
			FailedIds=append(FailedIds,reservation.Id)
			continue
		}
	}
	if len(FailedIds)>0{
		return errorx.WrapError(errorx.UpdateErr,errors.New("部分预约记录状态更新失败"))
	}
	return nil
}
