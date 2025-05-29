package handler

import (
	"context"
	"net/http"

	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func reserveSeatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReserveSeatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		debug := r.Header.Get("DEBUG_MODE")
		ctx := context.WithValue(r.Context(), "DEBUG_MODE", debug)

		l := logic.NewReserveSeatLogic(ctx, svcCtx)
		resp, err := l.ReserveSeat(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
