package controllers

import (
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

}
