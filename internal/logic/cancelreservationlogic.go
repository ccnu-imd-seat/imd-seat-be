package logic

import (
	"context"
	"errors"
	"strconv"
	"time"

	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/contextx"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
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
	debug := l.ctx.Value("DEBUG_MODE")

	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}

	ID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return nil, errorx.WrapError(errorx.DefaultErr, errors.New("ID转换失败"))
	}

	//获取预约信息
	ReservationInfro, err := l.svcCtx.ReservationModel.FindOne(l.ctx, ID)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	} else if ReservationInfro.StudentId != studentID || ReservationInfro.Status != types.BookedStatus {
		return nil, errorx.WrapError(errorx.ViolateErr, errors.New("该预约不属于当前用户或状态不符合取消要求"))
	}

	//检验是否合规
	if !l.CheckCancelRule(ReservationInfro) && debug != "1" {
		return nil, errorx.WrapError(errorx.ViolateErr, errors.New("取消请求不合规"))
	}

	//更改预约状态
	err = l.svcCtx.ReservationModel.UpdateReservstionMessage(l.ctx, ID, types.CancelledStatus)
	if err != nil {
		return nil, errorx.WrapError(errorx.UpdateErr, err)
	}
	
	//释放座位，改状态为可预约
	err = l.svcCtx.SeatModel.ChangeSeatStatusByType(l.ctx, ReservationInfro.Date, types.AvaliableStatus, ReservationInfro.Seat, ReservationInfro.Type)
	if err != nil {
		return nil, errorx.WrapError(errorx.UpdateErr, err)
	}
	resp = &types.GeneralRes{
		Base: response.Success(),
	}
	return resp, nil
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
