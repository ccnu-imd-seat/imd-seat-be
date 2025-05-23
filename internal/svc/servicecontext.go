package svc

import (
	"net/http"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/middleware"
	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/ijwt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	JWTHandler       ijwt.JWTHandler
	AuthMiddleware   func(handlerFunc http.HandlerFunc) http.HandlerFunc
	SeatModel        model.SeatModel
	ReservationModel model.ReservationModel
	RoomModel        model.RoomModel
}

func NewServiceContext(c config.Config, conn sqlx.SqlConn) *ServiceContext {
	JWTHandler := ijwt.NewJWTHandler(c.Auth.AccessSecret)
	AuthMiddleware := middleware.NewAuthMiddleware(c, JWTHandler).AuthHandle

	return &ServiceContext{
		Config:           c,
		JWTHandler:       JWTHandler,
		SeatModel:        model.NewSeatModel(conn),
		ReservationModel: model.NewReservationModel(conn),
		AuthMiddleware:   AuthMiddleware,
		RoomModel:        model.NewRoomModel(conn, cache.ClusterConf{}),
	}
}
