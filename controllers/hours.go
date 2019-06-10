package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/common"
	"github.com/szdx4/attsys-server/utils/database"
)

// HoursShow 获取工时记录
func HoursShow(c *gin.Context) {
	// 初始化条件查询模型
	hours := []models.Hours{}
	db := database.Connector.Preload("User").Order("created_at DESC").Joins("LEFT JOIN users ON hours.user_id = users.id")

	// 获取认证用户信息
	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 从 URL 中获取用户 ID 并验证权限
	if userID, isExist := c.GetQuery("user_id"); isExist {
		userID, _ := strconv.Atoi(userID)

		if role == "user" && userID != authID {
			response.Unauthorized(c, "You cannot get others information")
			c.Abort()
			return
		}

		db = db.Where("user_id = ?", userID)
	} else if role == "user" {
		response.Unauthorized(c, "You cannot get others information")
		c.Abort()
		return
	}

	// 检测 start_at 格式
	if startAt, isExist := c.GetQuery("start_at"); isExist {
		startAt, err := common.ParseTime(startAt)
		if err != nil {
			response.BadRequest(c, "Invalid start_at format")
			c.Abort()
			return
		}
		db = db.Where("date >= ?", startAt)
	}

	// 检测 end_at 格式
	if endAt, isExist := c.GetQuery("end_at"); isExist {
		endAt, err := common.ParseTime(endAt)
		if err != nil {
			response.BadRequest(c, "Invalid start_at format")
			c.Abort()
			return
		}
		db = db.Where("date <= ?", endAt)
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

	// 部门主管只能获取本部门的工时记录
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		db = db.Where("users.department_id = ?", manager.DepartmentID)
	}

	// 执行查询
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&hours)
	db.Model(&models.Hours{}).Count(&total)

	// 判断当前页是否为空
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.HoursShow(c, total, page, hours)
}
