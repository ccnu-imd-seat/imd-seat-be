package svc

import (
	"net/http"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/middleware"
	"imd-seat-be/internal/model"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	AuthMiddleware   func(handlerFunc http.HandlerFunc) http.HandlerFunc
	SeatModel        model.SeatModel
	ReservationModel model.ReservationModel
	RoomModel        model.RoomModel
}

func NewServiceContext(c config.Config, conn sqlx.SqlConn) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		AuthMiddleware :middleware.NewAuthMiddleware(c).AuthHandle,
		SeatModel:        model.NewSeatModel(conn, cache.ClusterConf{}), // 缓存配置为空，暂不启用缓存逻辑
		ReservationModel: model.NewReservationModel(conn, cache.ClusterConf{}),
		RoomModel:        model.NewRoomModel(conn, cache.ClusterConf{}),
	}
}
