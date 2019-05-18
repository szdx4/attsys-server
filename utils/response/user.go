package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserAuth 用户认证响应
func UserAuth(c *gin.Context, token string) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  token,
	})
}
