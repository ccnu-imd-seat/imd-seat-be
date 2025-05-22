package handler

import (
	"net/http"

	"imd-seat-be/internal/logic"
	"imd-seat-be/internal/pkg/ijwt"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			err = svcCtx.JWTHandler.SetJWTToken(w, ijwt.ClaimParams{
				StudentId: req.Username,
				Password:  req.Password,
			})
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
			}
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
