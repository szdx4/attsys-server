package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
)

// SignGetQrcode 获取二维码
func SignGetQrcode(c *gin.Context) {
	qrcode := models.Qrcode{
		ExpiredAt: time.Now().Add(60 * time.Second),
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
	database.Connector.Find(&sign, signID)

	if sign.ID == 0 {
		response.NotFound(c, "Sign not found")
		c.Abort()
		return
	}

	sign.EndAt = time.Now()
	database.Connector.Save(&sign)

	shift := models.Shift{}
	database.Connector.First(&shift, sign.ShiftID)

	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	shift.Status = "off"
	database.Connector.Save(&shift)

	timeDiff := uint(sign.EndAt.Sub(sign.StartAt).Hours())

	user := shift.User
	user.Hours += timeDiff
	database.Connector.Save(&user)

	hours := models.Hours{
		UserID: user.ID,
		Date:   sign.EndAt,
		Hours:  timeDiff,
	}
	database.Connector.Create(&hours)

	response.SignOff(c)
}
