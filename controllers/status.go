package controllers

import (
	"github.com/szdx4/attsys-server/utils/common"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
)

// StatusUser 用户相关数据
func StatusUser(c *gin.Context) {
	// 初始化条件查询
	db := database.Connector

	// 根据用户权限筛选不同数据
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		db.Where("department_id = ?", manager.DepartmentID)
	}

	// 初始化结果集
	usersCount := 0

	// 执行查询
	db.Model(&models.User{}).Count(&usersCount)

	response.StatusUser(c, usersCount)
}

// StatusSign 签到相关数据
func StatusSign(c *gin.Context) {
	// 初始化数据库操作对象
	db := database.Connector

	// 获取认证用户信息
	authID, _ := c.Get("user_id")
	authUser := models.User{}
	db.First(&authUser, authID)

	// 查询已签到用户
	signedCount := 0
	signedDb := db.Joins("LEFT JOIN users ON users.id = shifts.user_id").Where("status = 'on'")
	if authUser.Role == "manager" {
		signedDb = signedDb.Where("users.department_id = ?", authUser.DepartmentID)
	}
	signedDb.Model(&models.Shift{}).Count(&signedCount)

	// 查询迟到用户
	latedCount := 0

	// 还没有签到
	latedSignedCount := 0
	latedDb := db.Select("user_id, COUNT(*)").Where("status = 'no' AND start_at < ?", time.Now())
	latedDb.Model(&models.Shift{}).Group("user_id").Count(&latedSignedCount)

	// 已经签到
	latedUnsignedCount := 0
	latedDb = db.Select("user_id, COUNT(*)").Joins("LEFT JOIN signs ON signs.shift_id = shifts.id").Where("signs.start_at > shifts.start_at")
	latedDb.Model(&models.Shift{}).Group("user_id").Count(&latedSignedCount)

	// 加起来
	latedCount = latedSignedCount + latedUnsignedCount

	// 查询请假用户
	leavedCount := 0
	leavedDb := db.Select("user_id, COUNT(*)").Where("status = 'pass' AND start_at < ? AND end_at > ?", time.Now(), time.Now())
	leavedDb.Model(&models.Leave{}).Count(&leavedCount)

	// 发送响应
	response.StatusSign(c, signedCount, latedCount, leavedCount)
}

// StatusHour 获取用户工作时间和加班时间
func StatusHour(c *gin.Context) {
	// 初始化条件查询模型
	shifts := []models.Shift{}
	overtimes := []models.Overtime{}

	dbShift := database.Connector
	dbOvertime := database.Connector

	// 从URL中获取用户ID
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "user ID not valid")
		c.Abort()
		return
	}

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 验证用户
	if role == "user" && userID != authID {
		response.Unauthorized(c, "You cannot get other information")
		c.Abort()
		return
	}

	// 验证部门主管
	if role == "manager" && userID != authID {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		aim := models.User{}
		database.Connector.First(&aim, userID)
		if manager.DepartmentID != aim.DepartmentID {
			response.Unauthorized(c, "You cannot get other department information")
			c.Abort()
			return
		}
	}

	// 查找对应 user
	dbShift = dbShift.Where("user_id = ?", userID)
	dbOvertime = dbOvertime.Where("user_id = ?", userID)

	// 检测 start_at
	if startAt, isExist := c.GetQuery("start_at"); isExist {
		startAt, err := common.ParseTime(startAt)
		if err != nil {
			response.BadRequest(c, "Invalid start_at format")
			c.Abort()
			return
		}
		dbShift = dbShift.Where("start_at >= ?", startAt)
		dbOvertime = dbOvertime.Where("start_at >= ?", startAt)
	} else {
		response.BadRequest(c, "start_at absents")
		c.Abort()
		return
	}

	// 检测 end_at 格式
	if endAt, isExist := c.GetQuery("end_at"); isExist {
		endAt, err := common.ParseTime(endAt)
		if err != nil {
			response.BadRequest(c, "Invalid end_at format")
			c.Abort()
			return
		}
		dbShift = dbShift.Where("end_at <= ?", endAt)
		dbOvertime = dbOvertime.Where("end_at <= ?", endAt)
	} else {
		response.BadRequest(c, "end_at absents")
		c.Abort()
		return
	}

	// 获得查询模型列表
	dbShift = dbShift.Where("status = ?", "off")
	dbOvertime = dbOvertime.Where("status = ?", "pass")
	dbShift.Find(&shifts)
	dbOvertime.Find(&overtimes)

	var shiftHour int
	var overtimeHour int

	// 统计 shift 中排班和加班
	for _, shift := range shifts {
		if shift.Type == "normal" {
			shiftHour += int(shift.EndAt.Sub(shift.StartAt).Hours())
		} else {
			overtimeHour += int(shift.EndAt.Sub(shift.StartAt).Hours())
		}
	}

	// 统计 overtime 中加班
	for _, overtime := range overtimes {
		overtimeHour += int(overtime.EndAt.Sub(overtime.StartAt).Hours())
	}

	// 发送响应
	response.StatusHour(c, shiftHour, overtimeHour)

}
