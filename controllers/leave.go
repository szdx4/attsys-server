package controllers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
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

	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 构造并存入数据库
	startAt, err := config.StrToTime(req.StartAt)
	if err != nil {
		response.BadRequest(c, errors.New("start_at not valid").Error())
		c.Abort()
		return
	}
	endAt, err := config.StrToTime(req.EndAt)
	if err != nil {
		response.BadRequest(c, errors.New("end_at not valid").Error())
		c.Abort()
		return
	}

	userID, _ := strconv.Atoi(c.Param("id"))
	authID, _ := c.Get("user_id")

	if userID != authID {
		response.Unauthorized(c, "You can only apply leave for yourself")
		c.Abort()
		return
	}

	user := models.User{}
	database.Connector.First(&user, userID)

	leave := models.Leave{
		UserID:  uint(userID),
		StartAt: startAt,
		EndAt:   endAt,
		Remark:  req.Remark,
		Status:  "wait",
	}

	database.Connector.Create(&leave)
	if leave.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	managerID := user.Department.ManagerID
	message.Send(user.ID, managerID, "请假申请", "理由："+leave.Remark)

	response.LeaveCreate(c, leave.ID)
}

// LeaveShow 获取指定用户请假
func LeaveShow(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	total := 0

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

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

	leaves := []models.Leave{}
	db := database.Connector
	db = db.Where("user_id = ?", userID)
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&leaves)
	db.Model(&models.Leave{}).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.LeaveShow(c, total, page, leaves)
}

// LeaveList 请假列表
func LeaveList(c *gin.Context) {
	// 检测 page
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	total := 0

	leaves := []models.Leave{}
	db := database.Connector.Joins("LEFT JOIN users ON users.id = leaves.user_id")

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	if role == "manager" {
		// 用户只能获得本部门的请假列表
		manager := models.User{}
		database.Connector.First(&manager, authID)
		db = db.Where("users.department_id = ?", manager.DepartmentID)
	}

	db.Limit(perPage).Offset((page - 1) * perPage).Find(&leaves)
	db.Model(&models.Leave{}).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

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

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	leaveID, _ := strconv.Atoi(c.Param("id"))
	leave := models.Leave{}
	database.Connector.Preload("User").First(&leave, leaveID)
	if leave.ID == 0 {
		response.NotFound(c, "Leave not found")
		c.Abort()
		return
	}

	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		if manager.DepartmentID != leave.User.DepartmentID {
			response.Unauthorized(c, "You can only edit your department leave")
			c.Abort()
			return
		}
	}

	// 修改 leave 的 status
	leave.Status = req.Status
	database.Connector.Save(&leave)

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

	response.LeaveUpdate(c)
}

// LeaveDelete 销假
func LeaveDelete(c *gin.Context) {
	leaveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Leave ID invalid")
		c.Abort()
		return
	}

	leave := models.Leave{}
	database.Connector.First(&leave, leaveID)
	if leave.ID == 0 {
		response.NotFound(c, "Leave not found")
		c.Abort()
		return
	}

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

	// 销假之后将后续排班
	currentTime := time.Now()
	shifts := []models.Shift{}
	database.Connector.Where("start_at >= ? AND end_at <= ?", currentTime, leave.EndAt).Find(&shifts)
	for _, shift := range shifts {
		shift.Status = "no"
		database.Connector.Save(&shift)
	}

	leave.Status = "discarded"
	database.Connector.Save(&leave)

	response.LeaveDelete(c)
}
