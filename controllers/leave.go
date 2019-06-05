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
	leave := models.Leave{
		UserID:  uint(userID),
		StartAt: startAt,
		EndAt:   endAt,
		Remark:  req.Remark,
		Status:  "wait",
	}

	database.Connector.Create(&leave)
	if leave.ID < 1 {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	response.LeaveCreate(c, leave.ID)
}

// 获取指定用户请假
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

	leaves := []models.Leave{}
	db := database.Connector
	db = db.Where("user_id = ?", userID)
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&leaves)
	db.Model(&leaves).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.LeaveShow(c, total, page, leaves)
}

// 请假列表
func LeaveList(c *gin.Context) {
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
	database.Connector.Limit(perPage).Offset((page - 1) * perPage).Find(&leaves)
	database.Connector.Model(&leaves).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.LeaveList(c, total, page, leaves)
}

// 审批请假
func LeaveUpdate(c *gin.Context) {
	var req requests.LeaveUpdateRequest
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

	leaveID, _ := strconv.Atoi(c.Param("id"))
	leave := models.Leave{}
	database.Connector.Where("id = ?", leaveID).First(&leave)
	if leave.ID == 0 {
		response.NotFound(c, "shift not found")
		c.Abort()
		return
	}

	// 修改 leave 的 status
	leave.Status = req.Status
	database.Connector.Save(&leave)

	response.LeaveUpdate(c)
}
