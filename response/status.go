package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusUser 用户相关数据响应
func StatusUser(c *gin.Context, users int) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"users":  users,
	})
}

// StatusSign 用户相关数据响应
func StatusSign(c *gin.Context, signed, lated, leaved int) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"signed": signed,
		"lated":  lated,
		"leaved": leaved,
	})
}
