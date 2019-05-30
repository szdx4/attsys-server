package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
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
		"data":   user,
	})
}

// UserCreate 创建用户响应
func UserCreate(c *gin.Context, userID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"user_id": userID,
	})
}

// UserList 用户列表响应
func UserList(c *gin.Context, total, page int, users []models.User) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         users,
	})
}

// UserDelete 删除用户响应
func UserDelete(c *gin.Context, userID int) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"user_id": userID,
	})
}
