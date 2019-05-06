package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/database"
)

// DepartmentCreate 创建部门
func DepartmentCreate(c *gin.Context) {
	var req requests.CreateDepartmentRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "error",
		})
		return
	}

	department := models.Department{
		Name: req.Name,
	}
	database.Connector.Create(&department)

	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"resourseId": department.ID,
	})
}
