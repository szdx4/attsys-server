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
	list := []gin.H{}
	for _, shift := range shifts {
		list = append(list, gin.H{
			"id":         shift.ID,
			"user_id":    shift.User.ID,
			"user_name":  shift.User.Name,
			"start_at":   shift.StartAt,
			"end_at":     shift.EndAt,
			"type":       shift.Type,
			"status":     shift.Status,
			"created_at": shift.CreatedAt,
			"updated_at": shift.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         list,
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
