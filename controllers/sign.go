package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
)

// SignGetQrcode 获取二维码
func SignGetQrcode(c *gin.Context) {
	qrcode := models.Qrcode{
		ExpiredAt: time.Now().Add(time.Duration(config.App.QrcodeValidMinutes) * time.Minute),
	}
	qrcode.RandToken()
	database.Connector.Create(&qrcode)

	if qrcode.ID == 0 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	image, err := qrcode.Image()
	if err != nil {
		response.InternalServerError(c, "Qrcode generate error")
		c.Abort()
		return
	}

	response.SignGetQrcode(c, image, qrcode.ExpiredAt)
}

// SignWithQrcode 通过二维码签到
func SignWithQrcode(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	var req requests.SignWithQrcodeRequest
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

	// 获得登录者得 id
	authID, _ := c.Get("user_id")

	// 用户只能获取自己的签到情况
	if userID != authID {
		response.Unauthorized(c, "You can only sign for yourself")
		c.Abort()
		return
	}

	shift := models.Shift{}
	database.Connector.Where("status = 'no' AND user_id = ? AND end_at >= ?", userID, time.Now()).Order("start_at ASC").First(&shift)

	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	shift.Status = "on"
	database.Connector.Save(&shift)

	sign := models.Sign{
		ShiftID: shift.ID,
		StartAt: time.Now(),
		EndAt:   time.Now(),
	}
	database.Connector.Create(&sign)

	if sign.ID == 0 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	response.Sign(c, sign.ID)
}

// SignWithFace 通过人脸签到
func SignWithFace(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	var req requests.SignWithFaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	if err := req.Validate(userID); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	shift := models.Shift{}
	database.Connector.Where("status = 'no' AND user_id = ? AND end_at >= ?", userID, time.Now()).Order("start_at ASC").First(&shift)

	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	shift.Status = "on"
	database.Connector.Save(&shift)

	sign := models.Sign{
		ShiftID: shift.ID,
		StartAt: time.Now(),
		EndAt:   time.Now(),
	}
	database.Connector.Create(&sign)

	if sign.ID == 0 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	response.Sign(c, sign.ID)
}

// SignOff 签退
func SignOff(c *gin.Context) {
	signID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Sign ID invalid")
		c.Abort()
		return
	}

	sign := models.Sign{}
	database.Connector.Preload("Shift").Preload("User").First(&sign, signID)

	if sign.ID == 0 {
		response.NotFound(c, "Sign not found")
		c.Abort()
		return
	}

	diff := time.Now().Sub(sign.Shift.EndAt).Minutes()
	canOvertime := false

	if diff > float64(config.App.MinOvertimeMinutes) {
		sign.EndAt = sign.Shift.EndAt
		canOvertime = true
	} else {
		sign.EndAt = time.Now()
	}
	database.Connector.Save(&sign)

	sign.Shift.Status = "off"
	database.Connector.Save(&sign.Shift)

	timeDiff := uint(sign.EndAt.Sub(sign.StartAt).Hours())

	user := sign.Shift.User
	user.Hours += timeDiff
	database.Connector.Save(&user)

	hours := models.Hours{
		UserID: user.ID,
		Date:   sign.EndAt,
		Hours:  timeDiff,
	}
	database.Connector.Create(&hours)

	response.SignOff(c, canOvertime)
}

// SignStatus 获取用户签到状态
func SignStatus(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 用户只能获取自己的签到情况
	if role == "user" && userID != authID {
		response.Unauthorized(c, "You can only get your own sign status")
		c.Abort()
		return
	}

	// 部门主管只能获得本部门的员工签到情况
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		user := models.User{}
		database.Connector.First(&user, userID)
		if user.ID == 0 {
			response.NotFound(c, "User not found")
			c.Abort()
			return
		}
		if manager.DepartmentID != user.DepartmentID {
			response.Unauthorized(c, "You can only get your department sign status")
			c.Abort()
			return
		}
	}

	// 确认用户的存在性
	user := models.User{}
	database.Connector.First(&user, userID)
	if user.ID == 0 {
		response.NotFound(c, "User not found")
		c.Abort()
		return
	}

	// 查找用户的已经签到的排班
	shift := models.Shift{}
	database.Connector.Where("status = 'on'").First(&shift)
	if shift.ID == 0 {
		response.NoContent(c)
		c.Abort()
		return
	}

	sign := models.Sign{}
	database.Connector.Where("shift_id = ?", shift.ID).First(&sign)
	if sign.ID == 0 {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.SignStatus(c, sign.ID)
}
