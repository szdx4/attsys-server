package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
)

// LeaveCreate 申请请假响应
func LeaveCreate(c *gin.Context, leaveID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"user_id": leaveID,
	})
}

// LeaveShow 获取指定用户请假响应
func LeaveShow(c *gin.Context, total, page int, leaves []models.Leave) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         leaves,
	})
}

// LeaveList 请假列表响应
func LeaveList(c *gin.Context, total, page int, leaves []models.Leave) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         leaves,
	})
}

// LeaveUpdate 审批请假响应
func LeaveUpdate(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
	})
}

// LeaveDelete 销假响应
func LeaveDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
