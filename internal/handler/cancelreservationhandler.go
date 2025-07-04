package handler

import (
	"context"
	"net/http"

	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func cancelReservationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		debug := r.Header.Get("DEBUG_MODE")
		ctx := context.WithValue(r.Context(), "DEBUG_MODE", debug)

		var req types.CancelReservationReq
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCancelReservationLogic(ctx, svcCtx)
		resp, err := l.CancelReservation(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
