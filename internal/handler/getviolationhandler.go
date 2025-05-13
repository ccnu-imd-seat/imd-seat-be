package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
)

func getViolationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetViolationLogic(r.Context(), svcCtx)
		resp, err := l.GetViolation()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
