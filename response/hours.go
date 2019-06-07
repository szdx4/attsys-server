package response

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"net/http"
)

// HoursShow 获取工时记录响应
func HoursShow(c *gin.Context, total, page int, data []models.HourData) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         data,
	})
}
