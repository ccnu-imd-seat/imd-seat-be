package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/pkg/contextx"
	"imd-seat-be/internal/pkg/ijwt"
	"imd-seat-be/internal/pkg/response"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	Cfg config.Config
	r   ijwt.JWTHandler
}

func NewAuthMiddleware(cfg config.Config, r ijwt.JWTHandler) *AuthMiddleware {
	return &AuthMiddleware{Cfg: cfg,
		r: r}
}

func (m *AuthMiddleware) AuthHandle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")
		if AuthHeader == "" || !strings.HasPrefix(AuthHeader, "Bearer ") {
			err := fmt.Errorf("invalid authorization")
			code, body := response.ErrHandler(err)
			w.WriteHeader(code)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusOK, body)
			return
		}

		token := strings.TrimPrefix(AuthHeader, "Bearer ")
		claims, err := m.r.ParseToken(token)
		if err != nil {
			logx.Errorf("解析token失败: %v", err)
			code, body := response.ErrHandler(err)
			w.WriteHeader(code)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusOK, body)
			return
		}

		// 将学号信息写入context
		ctx := contextx.SetStudentID(r.Context(), claims.StudentId)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
