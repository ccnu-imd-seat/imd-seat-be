package task

import (
	"context"
	"errors"
	"log"
	"time"

	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

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
		log.Println("注册定时更新状态任务失败:", err)
	}
	//每月进行更新
	_, err = c.AddFunc("0 0 0 1 * *", func() {
		if err := RenewScore(ctx, svcCtx); err != nil {
			log.Println("更新信誉分失败:", err)
		} else {
			log.Println("更新信誉分成功")
		}
	})
	if err != nil {
		log.Println("注册定时更新信誉分失败:", err)
	}
	// 每天 0:05:00 清理过期座位信息
	_, err = c.AddFunc("0 5 0 * * *", func() {
		if err := CleanExpiredSeats(ctx, svcCtx); err != nil {
			log.Println("清理过期座位信息失败:", err)
		} else {
			log.Println("清理过期座位信息成功")
		}
	})
	if err != nil {
		log.Println("注册定时清理座位信息失败:", err)
	}
	c.Start()

	go func() {
		<-ctx.Done()
		c.Stop()
	}()
}

// 查找所有未签到的预约,更新状态为违约
func Violation(ctx context.Context, svcCtx *svc.ServiceContext) error {
	now := time.Now().Format("2006-01-02")
	parsedTime, err := time.Parse("2006-01-02", now)
	if err != nil {
		return errorx.WrapError(errorx.DefaultErr, errors.New("解析时间错误"))
	}
	Reservations, err := svcCtx.ReservationModel.GetReservationByStatus(ctx, parsedTime, types.EffectiveStatus)
	if err != nil {
		return errorx.WrapError(errorx.FetchErr, err)
	}

	for _, reservation := range Reservations {
		err := svcCtx.ReservationModel.UpdateReservstionMessage(ctx, reservation.Id, types.ViolatedStatus)
		if err != nil {
			log.Printf("更新id:%d的预约状态失败%v", reservation.Id, err)
			continue
		}
		//扣除信誉分
		err = ReduceScore(ctx, svcCtx, reservation.StudentId)
		if err != nil {
			log.Printf("更新用户:%s信誉分失败:%v", reservation.StudentId, err)
		}
	}
	return nil
}

// 扣除信誉分
func ReduceScore(ctx context.Context, svcCtx *svc.ServiceContext, StudentID string) error {
	score, err := svcCtx.UserModel.FindScoreByID(ctx, StudentID)
	if err != nil {
		return errorx.WrapError(errorx.FetchErr, err)
	}
	score = score - 100
	if score < 0 {
		score = 0
	}
	err = svcCtx.UserModel.UpdateScore(ctx, StudentID, score)
	if err != nil {
		return errorx.WrapError(errorx.UpdateErr, err)
	}
	return nil
}

// 每月恢复信誉分
func RenewScore(ctx context.Context, svcCtx *svc.ServiceContext) error {
	err := svcCtx.UserModel.RenewScore(ctx)
	if err != nil {
		return errorx.WrapError(errorx.UpdateErr, err)
	}
	return nil
}

// 删除今天之前的座位信息
func CleanExpiredSeats(ctx context.Context, svcCtx *svc.ServiceContext) error {
	// 获取当前日期
	today := time.Now().Format("2006-01-02")

	// 调用模型层方法删除早于今天的座位信息
	err := svcCtx.SeatModel.DeleteSeatsBeforeDate(ctx, today)
	if err != nil {
		log.Println("清理座位信息失败:", err)
		return errorx.WrapError(errorx.DeleteErr, err)
	}

	log.Println("过期座位信息清理成功")
	return nil
}
