package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
	"strconv"
)

// DepartmentCreate 创建部门
func DepartmentCreate(c *gin.Context) {
	var req requests.DepartmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Bad Request")
		c.Abort()
		return
	}
	//检验
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}
	//新建
	department := models.Department{
		Name:      req.Name,
		ManagerID: req.ManagerId,
	}
	database.Connector.Create(&department)

	if department.ID < 1 {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}
	//响应
	response.Created(c, department.ID)
}

// DepartmentList 获取部门列表
func DepartmentList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	departments := []models.Department{}
	total := 0
	database.Connector.Limit(perPage).Offset((page - 1) * perPage).Find(&departments)
	database.Connector.Model(&models.Department{}).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}
	response.DepartmentList(c, total, page, departments)
}
