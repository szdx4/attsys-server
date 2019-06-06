package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SignGetQrcode 获取二维码响应
func SignGetQrcode(c *gin.Context, image string, expiredAt time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"qrcode":     image,
			"expired_at": expiredAt,
		},
	})
}

// SignWithQrcode 二维码签到响应
func SignWithQrcode(c *gin.Context, signID uint) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"sign_id": signID,
	})
}
