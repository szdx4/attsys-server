package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LeaveCreate 申请请假响应
func LeaveCreate(c *gin.Context, leaveID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"user_id": leaveID,
	})
}
