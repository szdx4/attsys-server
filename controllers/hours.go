package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
)

// HoursShow 获取工时记录
func HoursShow(c *gin.Context) {
	hours := []models.Hours{}
	total := 0
	db := database.Connector
	// 检测user_id
	if userID, isExist := c.GetQuery("user_id"); isExist == true {
		userID, _ := strconv.Atoi(userID)
		db = db.Where("user = ?", userID)
	}
	// 检测start_at
	if startAt, isExist := c.GetQuery("start_at"); isExist == true {
		db = db.Where("date >= ?", startAt)
	}
	// 检测end_at
	if endAt, isExist := c.GetQuery("end_at"); isExist == true {
		db = db.Where("date <= ?", endAt)
	}
	// 检测page
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	db = db.Limit(perPage).Offset((page - 1) * perPage)

	if err := db.Find(&hours).Error; err != nil {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.HoursShow(c, total, page, hours)
}
