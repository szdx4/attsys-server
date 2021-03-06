package response

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"net/http"
)

// OvertimeCreate 申请加班响应
func OvertimeCreate(c *gin.Context, overtimeID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"shift_id": overtimeID,
	})
}

// OvertimeShow 获取指定用户加班响应
func OvertimeShow(c *gin.Context, total, page int, overtime []models.Overtime) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         overtime,
	})
}

// OvertimeList 加班申请列表响应
func OvertimeList(c *gin.Context, total, page int, overtime []models.Overtime) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         overtime,
	})
}

// OvertimeUpdate 审批加班响应
func OvertimeUpdate(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
	})
}
