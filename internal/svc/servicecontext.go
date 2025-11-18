package svc

import (
	"net/http"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/middleware"
	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/ijwt"
	"imd-seat-be/internal/pkg/image"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	ImageUploader       image.Uploader
	JWTHandler          ijwt.JWTHandler
	AuthMiddleware      func(handlerFunc http.HandlerFunc) http.HandlerFunc
	AuthMiddlewareAdmin func(handlerFunc http.HandlerFunc) http.HandlerFunc
	SeatModel           model.SeatModel
	ReservationModel    model.ReservationModel
	UserModel           model.UserModel
	RoomModel           model.RoomModel
	FeedbackModel       model.FeedbackModel
}

func NewServiceContext(c config.Config, conn sqlx.SqlConn) *ServiceContext {
	JWTHandler := ijwt.NewJWTHandler(c.Auth.AccessSecret)
	AuthMiddleware := middleware.NewAuthMiddleware(c, JWTHandler).AuthHandle
	AuthMiddlewareAdmin := middleware.NewAuthMiddleware(c, JWTHandler).AuthHandleAdmin

	return &ServiceContext{
		Config:              c,
		JWTHandler:          JWTHandler,
		SeatModel:           model.NewSeatModel(conn),
		ReservationModel:    model.NewReservationModel(conn),
		AuthMiddleware:      AuthMiddleware,
		AuthMiddlewareAdmin: AuthMiddlewareAdmin,
		UserModel:           model.NewUserModel(conn),
		RoomModel:           model.NewRoomModel(conn),
		FeedbackModel:       model.NewFeedbackModel(conn),
		ImageUploader:       image.NewQiniuUploader(c),
	}
}

