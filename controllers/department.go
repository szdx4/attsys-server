package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
)

// DepartmentCreate 创建部门
func DepartmentCreate(c *gin.Context) {
	var req requests.DepartmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Bad Request")
		c.Abort()
		return
	}
	// 检验
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}
	// 新建
	department := models.Department{
		Name: req.Name,
	}
	database.Connector.Create(&department)

	if department.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}
	// 响应
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

// DepartmentShow 获取指定部门信息
func DepartmentShow(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Department ID invalid")
		c.Abort()
		return
	}

	department := models.Department{}
	database.Connector.First(&department, departmentID)
	if department.ID == 0 {
		response.NotFound(c, "Department not found")
		c.Abort()
		return
	}

	role, _ := c.Get("user_role")
	if role != "master" {
		userID, _ := c.Get("user_id")
		user := models.User{}
		database.Connector.First(&user, userID)
		if user.DepartmentID != uint(departmentID) {
			response.Unauthorized(c, "You can only get your department information")
			c.Abort()
			return
		}
	}

	response.DepartmentShow(c, department)
}

// DepartmentUpdate 编辑部门
func DepartmentUpdate(c *gin.Context) {
	// 合法性检验
	var req requests.DepartmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}
	if err := req.Validate(c); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Department ID invalid")
		c.Abort()
		return
	}

	// 从数据库中查找
	department := models.Department{}
	database.Connector.First(&department, departmentID)
	if department.ID == 0 {
		response.NotFound(c, "Department not found")
		c.Abort()
		return
	}

	// 编辑部门的相应信息
	department.Name = req.Name
	database.Connector.Save(&department)

	response.DepartmentUpdate(c)
}

// DepartmentDelete 删除部门
func DepartmentDelete(c *gin.Context) {
	DepartmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Department ID invalid")
		c.Abort()
		return
	}

	department := models.Department{}
	database.Connector.Where("id = ?", DepartmentID).First(&department)

	if department.ID == 0 {
		response.NotFound(c, "Department not found")
		c.Abort()
		return
	}
	database.Connector.Delete(&department)

	response.DepartmentDelete(c)
}
