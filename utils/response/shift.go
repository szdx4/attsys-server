package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShiftCreate 添加排班响应
func ShiftCreate(c *gin.Context, shiftID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"shift_id": shiftID,
	})
}
