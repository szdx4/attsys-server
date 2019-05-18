package middleware

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/utils/response"
)

// Token 验证 Token 中间件
func Token(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		response.Unauthorized(c, "This action need authorized")
		c.Abort()
		return
	}

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

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		expiredAt, err := time.Parse(time.UnixDate, claims["expired_at"].(string))
		if err != nil {
			response.Unauthorized(c, "Auth token not valid 1")
			c.Abort()
			return
		}

		if time.Now().UTC().After(expiredAt) {
			response.Unauthorized(c, "Auth token expired")
			c.Abort()
			return
		}

		c.Set("user_id", claims["id"].(string))
		c.Set("user_role", claims["role"].(string))
		c.Next()
	} else {
		response.Unauthorized(c, "Auth token not valid")
		c.Abort()
		return
	}
}