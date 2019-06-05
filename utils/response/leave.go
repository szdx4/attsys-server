package response

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"net/http"
)

// LeaveCreate 申请请假响应
func LeaveCreate(c *gin.Context, leaveID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"user_id": leaveID,
	})
}

// LeaveShow获取指定用户请假响应
func LeaveShow(c *gin.Context, total, page int, leaves []models.Leave) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         leaves,
	})
}
