package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"
)

func getFeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FeedbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetFeedbackLogic(r.Context(), svcCtx)
		resp, err := l.GetFeedback(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
