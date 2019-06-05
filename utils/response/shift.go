package response

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"net/http"
)

// ShiftCreate 添加排班响应
func ShiftCreate(c *gin.Context, shiftID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"shift_id": shiftID,
	})
}

// ShiftList 排班列表响应
func ShiftList(c *gin.Context, total, page int, shifts []models.Shift) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         shifts,
	})
}
