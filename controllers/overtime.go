package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
	"strconv"
)

// OvertimeCreate 申请加班
func OvertimeCreate(c *gin.Context) {
	var req requests.OvertimeCreateRequest
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
	overtime := models.Overtime{
		UserID:  uint(userID),
		StartAt: startAt,
		EndAt:   endAt,
		Remark:  req.Remark,
		Status:  "wait",
	}
	database.Connector.Create(&overtime)
	if overtime.ID < 1 {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	response.OvertimeCreate(c, overtime.ID)
}

// OvertimeShow 获取指定用户加班
func OvertimeShow(c *gin.Context) {
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

	overtime := []models.Overtime{}
	db := database.Connector
	db = db.Where("user_id = ?", userID)
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&overtime)
	db.Model(&overtime).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.OvertimeShow(c, total, page, overtime)
}

// OvertimeList 加班申请列表
func OvertimeList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	total := 0

	overtime := []models.Overtime{}
	database.Connector.Limit(perPage).Offset((page - 1) * perPage).Find(&overtime)
	database.Connector.Model(&overtime).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.OvertimeList(c, total, page, overtime)
}
