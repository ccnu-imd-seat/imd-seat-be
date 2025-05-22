package ijwt

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTHandler struct {
}

// ExtractToken 从请求中提取并返回 JWT
func (r *JWTHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

// SetJWTToken 生成并设置用户的 JWT
func (r *JWTHandler) SetJWTToken(ctx *gin.Context, cp ClaimParams) error {
	uc := UserClaims{
		StudentId: cp.StudentId,
		Password:  cp.Password,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString(token)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
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
