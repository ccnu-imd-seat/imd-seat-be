package handler

import (
	"context"
	"net/http"

	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getAvailableDaysHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Type := r.URL.Query().Get("type")

		debug := r.Header.Get("DEBUG_MODE")
		ctx := context.WithValue(r.Context(), "DEBUG_MODE", debug)

		l := logic.NewGetAvailableDaysLogic(ctx, svcCtx)
		resp, err := l.GetAvailableDays(Type)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
