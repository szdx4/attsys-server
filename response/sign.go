package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
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

// Sign 签到响应
func Sign(c *gin.Context, signID uint, user models.User) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"sign_id": signID,
		"user":    user,
	})
}

// SignOff 签退响应
func SignOff(c *gin.Context, canOvertime bool) {
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"overtime": canOvertime,
	})
}

// SignStatus 签到状态响应
func SignStatus(c *gin.Context, signID uint, shift models.Shift) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"sign_id": signID,
		"shift":   shift,
	})
}
