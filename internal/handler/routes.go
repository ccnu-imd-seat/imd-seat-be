package handler

import (
	"net/http"
	"time"

	"imd-seat-be/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: loginHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(time.Minute),
	)

	authRoutes := []rest.Route{
		{
			Method:  http.MethodDelete,
			Path:    "/reservation/cancel/:id",
			Handler: cancelReservationHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/reservation/days",
			Handler: getAvailableDaysHandler(serverCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/reservation/reserve",
			Handler: reserveSeatHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/reservation/seats",
			Handler: getSeatInfoHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/mine/reservations",
			Handler: getMyReservationHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/mine/score",
			Handler: getScoreHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/reservation/rooms",
			Handler: getRoomsHandler(serverCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/feedback",
			Handler: submitFeedbackHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/feedback",
			Handler: getFeedbackHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/checkin",
			Handler: checkInHandler(serverCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/upload",
			Handler: UploadSeatCsvHandler(serverCtx),
		},
	}

	server.AddRoutes(
		// 进行 with 操作后返回的依然是[]rest.Route
		rest.WithMiddleware(serverCtx.AuthMiddleware, authRoutes...),
		rest.WithPrefix("/api/v1"),
	)

}
