package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SignQrcode 获取二维码响应
func SignQrcode(c *gin.Context, image string, expiredAt time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"qrcode":     image,
			"expired_at": expiredAt,
		},
	})
}
