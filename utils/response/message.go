package response

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"net/http"
)

// MessageShow 获取指定消息
func MessageShow(c *gin.Context, message models.Message) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   message,
	})
}
