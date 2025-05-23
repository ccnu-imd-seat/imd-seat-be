package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/pkg/ijwt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	Cfg config.Config
	r   *ijwt.JWTHandler
}

func NewAuthMiddleware(cfg config.Config, r *ijwt.JWTHandler) *AuthMiddleware {
	return &AuthMiddleware{Cfg: cfg,
		r: r}
}

func (m *AuthMiddleware) AuthHandle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")
		if AuthHeader == "" || !strings.HasPrefix(AuthHeader, "Bearer ") {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid authoriztion"))
			return
		}
		token := strings.TrimPrefix(AuthHeader, "Bearer ")
		claims, err := m.r.ParseToken(token)
		if err != nil {
			logx.Errorf("解析token失败:%v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//将学号信息写入context
		ctx := context.WithValue(r.Context(), "student_id", claims.StudentId)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
