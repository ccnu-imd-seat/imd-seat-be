package ijwt

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type JWTHandler interface {
	SetJWTToken(w http.ResponseWriter, cp ClaimParams) error
	ParseToken(tokenStr string) (*UserClaims, error)
}

type JWTHandlerImpl struct {
	Secret []byte
}

func NewJWTHandler(secret string) JWTHandler {
	return &JWTHandlerImpl{
		Secret: []byte(secret),
	}
}

// SetJWTToken 生成并设置 JWT 到响应头
func (r *JWTHandlerImpl) SetJWTToken(w http.ResponseWriter, cp ClaimParams) error {
	uc := UserClaims{
		StudentId: cp.StudentId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &uc)
	tokenStr, err := token.SignedString(r.Secret)
	if err != nil {
		return err
	}
	w.Header().Set("x-jwt-token", tokenStr)
	return nil
}

func (r *JWTHandlerImpl) ParseToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// 将接口类型断言成结构类型
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// UserClaims 定义了 JWT 中用户相关的声明
type UserClaims struct {
	jwt.RegisteredClaims
	StudentId string `json:"student_id"`
}

type ClaimParams struct {
	StudentId string `json:"student_id"`
}
