package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
)

// DepartmentCreate 创建部门
func DepartmentCreate(c *gin.Context) {
	var req requests.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Bad Request")
		return
	}

	department := models.Department{
		Name: req.Name,
	}
	database.Connector.Create(&department)

	response.Created(c, department.ID)
}
