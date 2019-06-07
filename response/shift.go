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

// ShiftDepartment 部门排班响应
func ShiftDepartment(c *gin.Context, shiftIds []uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"shift_ids": shiftIds,
	})
}

// ShiftUpdate 更新排班状态响应
func ShiftUpdate(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
	})
}

// ShiftDelete 删除排班响应
func ShiftDelete(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
	})
}