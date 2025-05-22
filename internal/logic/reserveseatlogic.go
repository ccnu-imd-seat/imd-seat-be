package logic

import (
	"context"
	"time"

	"imd-seat-be/internal/model"
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
	resp=&types.ReserveSeatRes{}
	student_id, ok := l.ctx.Value("student_id").(string)
	if !ok{
		resp.Base.Code=401
		resp.Base.Message="解析学号失败"
		return
	}
	format:="2006-01-02"
	t,err:=time.Parse(format,req.Date)
	ReservationInfro :=model.Reservation{
		StudentId: student_id,
		Type: req.Type,
		Date:t,
		Room:req.Room,
		Seat: req.SeatID,
		Status: types.BookedStatus,
	}
	re,err:=l.svcCtx.ReservationModel.Insert(l.ctx,&ReservationInfro)
	if err!=nil{
		resp.Base.Code=500
		resp.Base.Message="预约失败"
		return
	}
	id,_:=re.LastInsertId()
	resp.Base.Code=200
	resp.Base.Message="预约成功"
	resp.Data=types.ReservationData{
		RoomID:req.Room,
		Seat: req.SeatID,
		Date:req.Date,
		Type: req.Type,
		ReservationID: int(id),
	}
	return
}
