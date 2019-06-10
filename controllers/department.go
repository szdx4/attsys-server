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

	// 检验提交数据的合法性
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 创建新部门模型并保存到数据库
	department := models.Department{
		Name: req.Name,
	}
	database.Connector.Create(&department)

	// 判断是否创建成功
	if department.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	// 发送响应
	response.Created(c, department.ID)
}

// DepartmentList 获取部门列表
func DepartmentList(c *gin.Context) {
	// 处理分页
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage

	// 按分页查询对应的数据
	departments := []models.Department{}
	total := 0
	database.Connector.Limit(perPage).Offset((page - 1) * perPage).Find(&departments)
	database.Connector.Model(&models.Department{}).Count(&total)

	// 判断当前页是否为空
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.DepartmentList(c, total, page, departments)
}

// DepartmentShow 获取指定部门信息
func DepartmentShow(c *gin.Context) {
	// 从 URL 中获得部门 ID
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Department ID invalid")
		c.Abort()
		return
	}

	// 数据库中查询部门
	department := models.Department{}
	database.Connector.First(&department, departmentID)
	if department.ID == 0 {
		response.NotFound(c, "Department not found")
		c.Abort()
		return
	}

	// 判断当前用户权限
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

	// 发送响应
	response.DepartmentShow(c, department)
}

// DepartmentUpdate 编辑部门
func DepartmentUpdate(c *gin.Context) {
	var req requests.DepartmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 检验提交数据的合法性
	if err := req.Validate(c); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 从 URL 中获取部门 ID
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Department ID invalid")
		c.Abort()
		return
	}

	// 从数据库中查找部门
	department := models.Department{}
	database.Connector.First(&department, departmentID)
	if department.ID == 0 {
		response.NotFound(c, "Department not found")
		c.Abort()
		return
	}

	// 编辑部门的相应信息并保存
	department.Name = req.Name
	database.Connector.Save(&department)

	// 发送响应
	response.DepartmentUpdate(c)
}

// DepartmentDelete 删除部门
func DepartmentDelete(c *gin.Context) {
	// 从 URL 中获取部门 ID
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Department ID invalid")
		c.Abort()
		return
	}

	// 从数据库中查找部门
	department := models.Department{}
	database.Connector.First(&department, departmentID)

	// 判断部门是否存在
	if department.ID == 0 {
		response.NotFound(c, "Department not found")
		c.Abort()
		return
	}

	// 执行删除操作
	database.Connector.Delete(&department)

	// 发送响应
	response.DepartmentDelete(c)
}
