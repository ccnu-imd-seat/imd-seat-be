package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"
)

func checkInHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckIn
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCheckInLogic(r.Context(), svcCtx)
		resp, err := l.CheckIn(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
