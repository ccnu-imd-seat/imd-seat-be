package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
)

func getRoomsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetRoomsLogic(r.Context(), svcCtx)
		resp, err := l.GetRooms()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
