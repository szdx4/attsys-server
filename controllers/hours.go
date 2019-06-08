package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
)

// HoursShow 获取工时记录
func HoursShow(c *gin.Context) {
	hours := []models.Hours{}
	db := database.Connector.Joins("LEFT JOIN users ON hours.user_id = users.id")

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 检测 user_id
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

	// 检测 start_at
	if startAt, isExist := c.GetQuery("start_at"); isExist {
		db = db.Where("date >= ?", startAt)
	}

	// 检测 end_at
	if endAt, isExist := c.GetQuery("end_at"); isExist {
		db = db.Where("date <= ?", endAt)
	}

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

	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		db = db.Where("users.department_id = ?", manager.DepartmentID)
	}

	db.Limit(perPage).Offset((page - 1) * perPage).Find(&hours)
	db.Model(&models.Hours{}).Count(&total)
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 构造 data 响应
	datas := []models.HourData{}
	for i := 0; i < len(hours); i++ {
		data := models.HourData{
			ID:       hours[i].ID,
			UserID:   hours[i].ID,
			UserName: hours[i].User.Name,
			Date:     hours[i].Date,
			Hours:    hours[i].Hours,
		}
		datas = append(datas, data)
	}

	response.HoursShow(c, total, page, datas)
}
