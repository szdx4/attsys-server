package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// OvertimeCreate 申请加班响应
func OvertimeCreate(c *gin.Context, overtimeID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"shift_id": overtimeID,
	})
}
