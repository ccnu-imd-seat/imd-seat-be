package handler

import (
	"net/http"
	"time"

	"imd-seat-be/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	anthRoutes := []rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/reservation/rooms",
			Handler: getRoomsHandler(serverCtx),
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
			Method:  http.MethodPost,
			Path:    "/new",
			Handler: NewSeatReqHandler(serverCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/upload",
			Handler: UploadSeatCsvHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/checkin",
			Handler: checkInHandler(serverCtx),
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
	}

	server.AddRoutes(
		rest.WithMiddleware(serverCtx.AuthMiddleware, anthRoutes...),
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(10*time.Minute),
	)

	adminRoutes := []rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/reservation/get_supreme_data",
			Handler: getSupremeDataHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/reservation/get_supreme_list",
			Handler: getSupremeListHandler(serverCtx),
		},
	}

	server.AddRoutes(
		rest.WithMiddleware(serverCtx.AuthMiddlewareAdmin, adminRoutes...),
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(10*time.Minute),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: loginHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(10*time.Minute),
	)

	server.AddRoute(
		rest.Route{
			Method:  http.MethodGet,
			Path:    "/download",
			Handler: download(serverCtx),
		},
		rest.WithTimeout(10*time.Minute),
	)
}
