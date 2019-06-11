package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/common"
	"github.com/szdx4/attsys-server/utils/database"
)

// ShiftCreate 添加排班
func ShiftCreate(c *gin.Context) {
	var req requests.ShiftCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 验证提交数据的合法性
	if err := req.Validate(c); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 解析 start_at 和 end_at 字段
	startAt, err := common.ParseTime(req.StartAt)
	if err != nil {
		response.BadRequest(c, "start_at not valid")
		c.Abort()
		return
	}
	endAt, err := common.ParseTime(req.EndAt)
	if err != nil {
		response.BadRequest(c, "end_at not valid")
		c.Abort()
		return
	}

	// 排班不能在现在之前
	if startAt.Before(time.Now()) {
		response.BadRequest(c, "You cannot arrange shift before now")
		c.Abort()
		return
	}

	// 获得 URL 中的用户 ID
	userID, _ := strconv.Atoi(c.Param("id"))

	// 创建排班
	shift := models.Shift{
		UserID:  uint(userID),
		StartAt: startAt,
		EndAt:   endAt,
		Type:    req.Type,
		Status:  "no",
	}
	database.Connector.Create(&shift)
	if shift.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	// 发送响应
	response.ShiftCreate(c, shift.ID)
}

// ShiftList 排班列表
func ShiftList(c *gin.Context) {
	// 初始化条件查询模型
	shifts := []models.Shift{}
	db := database.Connector.Preload("User").Joins("LEFT JOIN users ON shifts.user_id = users.id")

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 检测 user_id 参数
	if userID, isExit := c.GetQuery("user_id"); isExit {
		userID, _ := strconv.Atoi(userID)

		// 用户只能查询自己的排班
		if role == "user" && authID != userID {
			response.Unauthorized(c, "You can only get your information")
			c.Abort()
			return
		}

		db = db.Where("user_id = ?", userID)
	} else if role == "user" {
		// 用户不能得到列表
		response.Unauthorized(c, "You can only get your information")
		c.Abort()
		return
	}

	// 检测和解析 start_at 参数
	if startAt, isExit := c.GetQuery("start_at"); isExit {
		startAt, err := common.ParseTime(startAt)
		if err != nil {
			response.BadRequest(c, "invalid start_at format")
			c.Abort()
			return
		}
		db = db.Where("start_at >= ?", startAt)
	}

	// 检测和解析 end_at 参数
	if endAt, isExit := c.GetQuery("end_at"); isExit {
		endAt, err := common.ParseTime(endAt)
		if err != nil {
			response.BadRequest(c, "invalid end_at format")
			c.Abort()
			return
		}
		db = db.Where("end_at <= ?", endAt)
	}

	// 检测 department_id 参数
	if departmentID, isExit := c.GetQuery("department_id"); isExit {
		departmentID, _ := strconv.Atoi(departmentID)

		// 部门主管只能查看自己的部门
		if role == "manager" {
			manager := models.User{}
			database.Connector.First(&manager, authID)

			if manager.DepartmentID != uint(departmentID) {
				response.Unauthorized(c, "You can only get your department information")
				c.Abort()
				return
			}
		}

		db = db.Where("users.department_id = ?", departmentID)
	} else if role == "manager" {
		// 部门主管不能得到全体列表
		response.Unauthorized(c, "You can only get your department information")
		c.Abort()
		return
	}

	// 处理分页
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	total := 0

	// 执行查询
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&shifts)
	db.Model(&models.Shift{}).Count(&total)

	// 判断当前页是否为空
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.ShiftList(c, total, page, shifts)
}

// ShiftDepartment 部门排班
func ShiftDepartment(c *gin.Context) {
	// 检测 department_id 参数
	departmentID, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		response.BadRequest(c, "Department ID not valid")
		c.Abort()
		return
	}

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	var req requests.ShiftDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 验证提交数据的合法性
	if err := req.Validate(departmentID, role.(string), authID.(int)); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 查找出所有部门员工
	users := []models.User{}
	database.Connector.Where("department_id = ?", departmentID).Find(&users)

	// 获得起始时间
	startAt, err := common.ParseTime(req.StartAt)
	if err != nil {
		response.BadRequest(c, "start_at not valid")
		c.Abort()
		return
	}
	endAt, err := common.ParseTime(req.EndAt)
	if err != nil {
		response.BadRequest(c, "end_at not valid")
		c.Abort()
		return
	}

	// 存入数据库
	var shiftIDs []uint
	for _, user := range users {
		shift := models.Shift{
			UserID:  user.ID,
			StartAt: startAt,
			EndAt:   endAt,
			Type:    req.Type,
			Status:  "no",
		}
		database.Connector.Create(&shift)
		if shift.ID == 0 {
			response.InternalServerError(c, "Database error")
			c.Abort()
			return
		}
		shiftIDs = append(shiftIDs, user.ID)
	}

	// 发送响应
	response.ShiftDepartment(c, shiftIDs)
}

// ShiftAll 全单位排班
func ShiftAll(c *gin.Context) {
	var req requests.ShiftAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 验证提交数据的合法性
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 查找出所有员工和部门主管
	users := []models.User{}
	database.Connector.Where("role <> ?", "master").Find(&users)

	// 获得起始时间
	startAt, err := common.ParseTime(req.StartAt)
	if err != nil {
		response.BadRequest(c, "start_at not valid")
		c.Abort()
		return
	}
	endAt, err := common.ParseTime(req.EndAt)
	if err != nil {
		response.BadRequest(c, "end_at not valid")
		c.Abort()
		return
	}

	// 存入数据库
	var shiftIDs []uint
	for _, user := range users {
		shift := models.Shift{
			UserID:  user.ID,
			StartAt: startAt,
			EndAt:   endAt,
			Type:    req.Type,
			Status:  "no",
		}
		database.Connector.Create(&shift)
		if shift.ID == 0 {
			response.InternalServerError(c, "Database error")
			c.Abort()
			return
		}
		shiftIDs = append(shiftIDs, user.ID)
	}

	// 发送响应
	response.ShiftAll(c, shiftIDs)
}

// ShiftDelete 删除排班
func ShiftDelete(c *gin.Context) {
	// 获取 URL 中的排班 ID
	shiftID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Shift ID invalid")
		c.Abort()
		return
	}

	// 在数据库中查找到对应的排班
	shift := models.Shift{}
	database.Connector.Preload("User").First(&shift, shiftID)
	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 判断用户权限
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		if manager.DepartmentID != shift.User.DepartmentID {
			response.Unauthorized(c, "You can only delete your department shifts")
			c.Abort()
			return
		}
	}

	// 执行删除操作
	database.Connector.Delete(&shift)

	// 发送响应
	response.ShiftDelete(c)
}
