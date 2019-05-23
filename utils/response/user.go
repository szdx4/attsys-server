package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
)

// UserAuth 用户认证响应
func UserAuth(c *gin.Context, userID uint, token string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"user_id": userID,
		"token":   token,
	})
}

// UserShow 用户资料响应
func UserShow(c *gin.Context, user models.User) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"role":       user.Role,
			"department": user.DepartmentID,
			"hours":      user.Hours,
		},
	})
}

// UserCreate 创建用户响应
func UserCreate(c *gin.Context, userID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"user_id": userID,
	})
}
