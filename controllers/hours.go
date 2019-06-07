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
	db := database.Connector
	// 检测 user_id
	if userID, isExist := c.GetQuery("user_id"); isExist == true {
		userID, _ := strconv.Atoi(userID)
		db = db.Where("user_id = ?", userID)
	}

	// 检测 start_at
	if startAt, isExist := c.GetQuery("start_at"); isExist == true {
		db = db.Where("date >= ?", startAt)
	}

	// 检测 end_at
	if endAt, isExist := c.GetQuery("end_at"); isExist == true {
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
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&hours)
	db.Model(&hours).Count(&total)
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 构造 data 响应
	datas := []models.HourData{}
	for i := 0; i < len(hours); i++ {
		user := models.User{}
		database.Connector.Where("id = ?", hours[i].UserID).First(&user)
		data := models.HourData{
			ID:       hours[i].ID,
			UserID:   hours[i].ID,
			UserName: user.Name,
			Date:     hours[i].Date,
			Hours:    hours[i].Hours,
		}
		datas = append(datas, data)
	}

	response.HoursShow(c, total, page, datas)
}
