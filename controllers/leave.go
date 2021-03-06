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
	"github.com/szdx4/attsys-server/utils/message"
)

// LeaveCreate 申请请假
func LeaveCreate(c *gin.Context) {
	var req requests.LeaveCreateRequest
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

	// 解析 start_at 和 end_at 参数
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

	// 判断时间段的合法性
	if startAt.Before(time.Now()) {
		response.BadRequest(c, "You cannot apply leave before now")
		c.Abort()
		return
	}

	// 获取 URL 中的用户 ID 和认证用户 ID
	userID, _ := strconv.Atoi(c.Param("id"))
	authID, _ := c.Get("user_id")

	// 判断用户是否为自己请假
	if userID != authID {
		response.Unauthorized(c, "You can only apply leave for yourself")
		c.Abort()
		return
	}

	// 获取用户信息
	user := models.User{}
	database.Connector.First(&user, userID)

	// 建立新的请假模型
	leave := models.Leave{
		UserID:  uint(userID),
		StartAt: startAt,
		EndAt:   endAt,
		Remark:  req.Remark,
		Status:  "wait",
	}

	// 在数据库中插入数据
	database.Connector.Create(&leave)
	if leave.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	// 获取部门主管信息并发送信息
	manager := common.DepartmentManager(user.Department)
	if manager != nil {
		message.Send(user.ID, manager.ID, "请假申请", "理由："+leave.Remark)
	}

	// 发送响应
	response.LeaveCreate(c, leave.ID)
}

// LeaveShow 获取指定用户请假
func LeaveShow(c *gin.Context) {
	// 获取 URL 中的用户 ID
	userID, _ := strconv.Atoi(c.Param("id"))

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

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 判断用户权限
	if role == "user" && authID != userID {
		response.Unauthorized(c, "You can only get your leave")
		c.Abort()
		return
	}

	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		user := models.User{}
		database.Connector.First(&user, userID)
		if manager.DepartmentID != user.DepartmentID {
			response.Unauthorized(c, "You can only get your department leave")
			c.Abort()
			return
		}
	}

	// 查询请假信息
	leaves := []models.Leave{}
	db := database.Connector.Preload("User")
	db = db.Where("user_id = ?", userID)
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&leaves)
	db.Model(&models.Leave{}).Count(&total)

	// 判断当前页是否为空
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.LeaveShow(c, total, page, leaves)
}

// LeaveList 请假列表
func LeaveList(c *gin.Context) {
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

	// 初始化条件查询模型
	leaves := []models.Leave{}
	db := database.Connector.Preload("User").Joins("LEFT JOIN users ON users.id = leaves.user_id")

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 权限验证
	if role == "manager" {
		// 部门主管只能获得本部门的请假列表
		manager := models.User{}
		database.Connector.First(&manager, authID)
		db = db.Where("users.department_id = ?", manager.DepartmentID)
	}

	// 执行查询
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&leaves)
	db.Model(&models.Leave{}).Count(&total)

	// 判断当前页是否为空
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.LeaveList(c, total, page, leaves)
}

// LeaveUpdate 审批请假
func LeaveUpdate(c *gin.Context) {
	var req requests.LeaveUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 验证提交数据的合法性
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 查询请假信息
	leaveID, _ := strconv.Atoi(c.Param("id"))
	leave := models.Leave{}
	database.Connector.Preload("User").First(&leave, leaveID)
	if leave.ID == 0 {
		response.NotFound(c, "Leave not found")
		c.Abort()
		return
	}

	// 权限验证
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		if manager.DepartmentID != leave.User.DepartmentID {
			response.Unauthorized(c, "You can only edit your department leave")
			c.Abort()
			return
		}
	}

	// 修改请假申请的状态
	leave.Status = req.Status
	database.Connector.Save(&leave)

	// 检测状态，修改对应的排班状态并发送信息
	if leave.Status == "pass" {
		shifts := []models.Shift{}
		database.Connector.Where("start_at >= ? AND end_at <= ?", leave.StartAt, leave.EndAt).Find(&shifts)
		for _, shift := range shifts {
			shift.Status = "leave"
			database.Connector.Save(&shift)
		}
		message.Send(uint(authID.(int)), leave.UserID, "请假审批结果", "请假审批通过")
	} else {
		message.Send(uint(authID.(int)), leave.UserID, "请假审批结果", "请假审批未通过")
	}

	// 发送响应
	response.LeaveUpdate(c)
}

// LeaveDelete 销假
func LeaveDelete(c *gin.Context) {
	// 获取 URL 中的请假 ID
	leaveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Leave ID invalid")
		c.Abort()
		return
	}

	// 从数据库中查询请假
	leave := models.Leave{}
	database.Connector.First(&leave, leaveID)
	if leave.ID == 0 {
		response.NotFound(c, "Leave not found")
		c.Abort()
		return
	}

	// 获取认证用户信息
	authID, _ := c.Get("user_id")

	// 只能给自己销假
	if leave.UserID != uint(authID.(int)) {
		response.Unauthorized(c, "You can only cancel your own leave")
		c.Abort()
		return
	}

	// 只能给审核通过的请假销假
	if leave.Status != "pass" {
		response.BadRequest(c, "Leave is not passed")
		c.Abort()
		return
	}

	// 销假之后将后续排班状态设置为正常
	currentTime := time.Now()
	shifts := []models.Shift{}
	database.Connector.Where("start_at >= ? AND end_at <= ?", currentTime, leave.EndAt).Find(&shifts)
	for _, shift := range shifts {
		shift.Status = "no"
		database.Connector.Save(&shift)
	}

	// 更改请假申请状态
	leave.Status = "discarded"
	database.Connector.Save(&leave)

	// 发送响应
	response.LeaveDelete(c)
}
