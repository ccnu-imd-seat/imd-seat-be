package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"imd-seat-be/internal/config"
	"imd-seat-be/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	Cfg config.Config
}

func NewAuthMiddleware(cfg config.Config) *AuthMiddleware {
	return &AuthMiddleware{Cfg: cfg}
}

func (m *AuthMiddleware) AuthHandle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")
		if AuthHeader == "" || !strings.HasPrefix(AuthHeader, "Bearer ") {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid authoriztion"))
			return
		}
		token := strings.TrimPrefix(AuthHeader, "Bearer ")
		claims, err := utils.ParseToken(m.Cfg.Auth.AccessSecret, token)
		if err != nil {
			logx.Errorf("解析token失败:%v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//将学号信息写入context
		ctx := context.WithValue(r.Context(), "student_id", claims.StudentID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
