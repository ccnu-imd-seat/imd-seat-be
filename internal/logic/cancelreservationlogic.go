package logic

import (
	"context"
	"strconv"
	"time"

	"imd-seat-be/internal/model"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelReservationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelReservationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelReservationLogic {
	return &CancelReservationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelReservationLogic) CancelReservation(req *types.CancelReservationReq) (resp *types.GeneralRes, err error) {
	resp = &types.GeneralRes{}
	ID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Base.Code = 500
		resp.Base.Message = "ID 转换失败"
		return
	}
	//获取预约信息
	ReservationInfro, err := l.svcCtx.ReservationModel.FindOne(l.ctx, ID)
	if err != nil {
		resp.Base.Code = 500
		resp.Base.Message = "查询预约记录失败"
		return
	}
	if !l.CheckCancelRule(ReservationInfro) {
		resp.Base.Code = 400
		resp.Base.Message = "取消请求不合规"
		return
	}
	err = l.svcCtx.ReservationModel.Delete(l.ctx, ID)
	if err != nil {
		resp.Base.Code = 500
		resp.Base.Message = "取消预约失败"
		return
	}
	resp.Base.Code = 200
	resp.Base.Message = "取消预约成功"
	return
}

// 检验取消时间是否符合规则
func (l *CancelReservationLogic) CheckCancelRule(ReservationInfro *model.Reservation) bool {
	pre := ReservationInfro.Date.AddDate(0, 0, -1)
	target := time.Date(
		pre.Year(), pre.Month(), pre.Day(), 18, 0, 0, 0, pre.Location(),
	)
	now := time.Now().In(target.Location()).Truncate(time.Second)
	if ReservationInfro.Type == "week" || !now.Before(target) {
		return false
	}
	return true
}
