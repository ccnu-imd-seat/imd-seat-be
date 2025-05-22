package svc

import (
	"imd-seat-be/internal/config"
	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/ijwt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	JWTHandler       *ijwt.JWTHandler
	SeatModel        model.SeatModel
	ReservationModel model.ReservationModel
}

func NewServiceContext(c config.Config, conn sqlx.SqlConn) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		JWTHandler:       ijwt.NewJWTHandler(c.Auth.AccessSecret),
		SeatModel:        model.NewSeatModel(conn),
		ReservationModel: model.NewReservationModel(conn),
	}
}
