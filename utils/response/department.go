package response

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"net/http"
)

// DeparmentList 部门列表响应
func DepartmentList(c *gin.Context, total, page int, departments []models.Department) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     config.App.ItemsPerPage,
		"data":         departments,
	})
}

// DepartmentShow 部门资料响应
func DepartmentShow(c *gin.Context, department models.Department) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   department,
	})
}
