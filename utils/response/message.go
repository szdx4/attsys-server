package response

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
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

// MessageList 获取信息列表
func MessageList(c *gin.Context, total, page int, messages []models.Message) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         messages,
	})
}
