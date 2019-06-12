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

// StatusHour 获取用户时间和加班时间
func StatusHour(c *gin.Context, shiftHour int, overtimeHour int) {
	c.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"shift_hour":    shiftHour,
		"overtime_hour": overtimeHour,
	})
}
