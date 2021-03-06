package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
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
func ShiftDepartment(c *gin.Context, shiftIDs []uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"shift_ids": shiftIDs,
	})
}

// ShiftDelete 删除排班响应
func ShiftDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

// ShiftAll 全单位排班响应
func ShiftAll(c *gin.Context, shiftIDs []uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"shift_ids": shiftIDs,
	})
}

// ShiftUpdate 修改排班响应
func ShiftUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
