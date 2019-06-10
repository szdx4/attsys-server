package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/response"
)

// Token 验证 Token 中间件
func Token(c *gin.Context) {
	// 获得 HTTP Headers 中 Authorization 的值
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(c, "This action need authorized")
		c.Abort()
		return
	}

	// 将 Authorization 分段
	auth := strings.Fields(authHeader)
	if len(auth) < 2 {
		response.Unauthorized(c, "This action need authorized")
		c.Abort()
		return
	}

	// 判断 Token 类型（仅支持 Bearer Token）
	if auth[0] != "Bearer" {
		response.BadRequest(c, "Token type unsupported")
		c.Abort()
		return
	}

	// 获得 Token
	tokenString := auth[1]

	// 验证 Token 的有效性
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.App.EncryptKey), nil
	})
	if err != nil {
		response.Unauthorized(c, "Auth token not valid")
		c.Abort()
		return
	}

	// 解析 Token 中编码的信息
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		expiredAt, err := time.Parse(time.RFC3339, claims["expired_at"].(string))
		if err != nil {
			response.Unauthorized(c, "Auth token not valid")
			c.Abort()
			return
		}

		if time.Now().UTC().After(expiredAt) {
			response.Unauthorized(c, "Auth token expired")
			c.Abort()
			return
		}

		c.Set("user_id", int(claims["id"].(float64)))
		c.Set("user_role", claims["role"].(string))
		c.Next()
	} else {
		response.Unauthorized(c, "Auth token not valid")
		c.Abort()
		return
	}
}
