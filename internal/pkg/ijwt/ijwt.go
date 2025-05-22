package ijwt

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTHandler struct {
	Secret []byte
}

func NewJWTHandler(secret string) *JWTHandler {
	return &JWTHandler{
		Secret: []byte(secret),
	}
}

// ExtractToken
func (r *JWTHandler) ExtractToken(headerValue string) string {
	if headerValue == "" {
		return ""
	}
	segs := strings.Split(headerValue, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

// SetJWTToken 生成并设置 JWT 到响应头
func (r *JWTHandler) SetJWTToken(w http.ResponseWriter, cp ClaimParams) error {
	uc := UserClaims{
		StudentId: cp.StudentId,
		Password:  cp.Password,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString(r.Secret)
	if err != nil {
		return err
	}
	w.Header().Set("x-jwt-token", tokenStr)
	return nil
}

// UserClaims 定义了 JWT 中用户相关的声明
type UserClaims struct {
	jwt.RegisteredClaims
	StudentId string
	Password  string
}

type ClaimParams struct {
	StudentId string
	Password  string
}
