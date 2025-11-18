package handler

import (
	"context"
	"net/http"

	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func checkInHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		debug := r.Header.Get("DEBUG_MODE")
		debug_day := r.Header.Get("DEBUG_DAY")
		ctx := context.WithValue(r.Context(), "DEBUG_MODE", debug)
		ctx = context.WithValue(ctx, "DEBUG_DAY", debug_day)

		var req types.CheckIn
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCheckInLogic(ctx, svcCtx)
		resp, err := l.CheckIn(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
