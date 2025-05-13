package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"
)

func submitFeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FeedbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSubmitFeedbackLogic(r.Context(), svcCtx)
		resp, err := l.SubmitFeedback(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
