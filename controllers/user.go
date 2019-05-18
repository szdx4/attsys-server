package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/response"
)

// UserAuth 用户认证
func UserAuth(c *gin.Context) {
	var req requests.UserAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Bad Request")
		return
	}

	user, err := req.Validate()
	if err != nil {
		response.Unauthorized(c, "Wrong username or password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID,
		"role":       user.Role,
		"expired_at": time.Now().UTC().Add(time.Hour * 2).Format(time.UnixDate),
	})

	tokenString, err := token.SignedString([]byte(config.App.EncryptKey))
	if err != nil {
		response.InternalServerError(c, "Internal Server Error")
		return
	}

	response.UserAuth(c, tokenString)
}

// UserShow 获取单个用户信息
func UserShow(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "test",
	})
}
