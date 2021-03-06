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
	// 创建一个二维码
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

	// 获得二维码图片的 dataURL
	image, err := qrcode.Image()
	if err != nil {
		response.InternalServerError(c, "Qrcode generate error")
		c.Abort()
		return
	}

	// 发送响应
	response.SignGetQrcode(c, image, qrcode.ExpiredAt)
}

// SignWithQrcode 通过二维码签到
func SignWithQrcode(c *gin.Context) {
	// 获取 URL 中的用户 ID
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

	// 验证提交数据的合法性
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 获取认证用户信息
	authID, _ := c.Get("user_id")

	// 用户只能获取自己的签到情况
	if userID != authID {
		response.Unauthorized(c, "You can only sign for yourself")
		c.Abort()
		return
	}

	// 验证是否重复签到
	shift := models.Shift{}
	database.Connector.Where("status = 'on' AND user_id = ?", userID).First(&shift)
	if shift.ID != 0 {
		response.BadRequest(c, "Shift had signed on")
		c.Abort()
		return
	}

	// 找到下一个未签到的排班
	database.Connector.Where("status = 'no' AND user_id = ? AND end_at >= ?", userID, time.Now()).Order("start_at ASC").First(&shift)
	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	// 将排班状态更改为已签到
	shift.Status = "on"
	database.Connector.Save(&shift)

	// 建立签到记录
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

	// 获取用户信息
	user := models.User{}
	database.Connector.First(&user, userID)

	// 发送响应
	response.Sign(c, sign.ID, user)
}

// SignWithFace 通过人脸签到
func SignWithFace(c *gin.Context) {
	var req requests.SignWithFaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 验证提交数据的合法性
	userID, err := req.Validate()
	if err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 获取用户信息
	user := models.User{}
	database.Connector.First(&user, userID)

	// 验证是否重复签到
	shift := models.Shift{}
	database.Connector.Where("status = 'on' AND user_id = ?", userID).First(&shift)
	if shift.ID != 0 {
		response.BadRequest(c, "Shift had signed on")
		c.Abort()
		return
	}

	// 找到下一个未签到的排班
	database.Connector.Where("status = 'no' AND user_id = ? AND end_at >= ?", userID, time.Now()).Order("start_at ASC").First(&shift)
	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	// 将排班状态更改为已签到
	shift.Status = "on"
	database.Connector.Save(&shift)

	// 建立签到记录
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

	// 发送响应
	response.Sign(c, sign.ID, user)
}

// SignOff 签退
func SignOff(c *gin.Context) {
	// 获取 URL 中的签到 ID
	signID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Sign ID invalid")
		c.Abort()
		return
	}

	// 查询签到记录
	sign := models.Sign{}
	database.Connector.Preload("Shift").Preload("Shift.User").First(&sign, signID)
	if sign.ID == 0 {
		response.NotFound(c, "Sign not found")
		c.Abort()
		return
	}

	// 计算能否申请加班
	diff := time.Now().Sub(sign.Shift.EndAt).Minutes()
	canOvertime := false
	if diff > float64(config.App.MinOvertimeMinutes) {
		sign.EndAt = sign.Shift.EndAt
		canOvertime = true
	} else {
		sign.EndAt = time.Now()
	}

	// 保存签到记录
	database.Connector.Save(&sign)

	// 判断是否重复签退
	if sign.Shift.Status == "off" {
		response.BadRequest(c, "Shift had signed off")
		c.Abort()
		return
	}

	// 更改排班状态
	sign.Shift.Status = "off"
	database.Connector.Save(&sign.Shift)

	// 计算工时
	timeDiff := uint(sign.Shift.EndAt.Sub(sign.StartAt).Hours())

	// 为用户加工时
	user := sign.Shift.User
	user.Hours += timeDiff
	database.Connector.Save(&user)

	// 创建工时记录
	hours := models.Hours{
		UserID: user.ID,
		Date:   sign.EndAt,
		Hours:  timeDiff,
	}
	database.Connector.Create(&hours)

	// 发送响应
	response.SignOff(c, canOvertime)
}

// SignStatus 获取用户签到状态
func SignStatus(c *gin.Context) {
	// 获取 URL 中的用户 ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	// 获取认证用户信息
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

	// 查找排班对应的签到记录
	sign := models.Sign{}
	database.Connector.Where("shift_id = ?", shift.ID).First(&sign)
	if sign.ID == 0 {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.SignStatus(c, sign.ID, shift)
}
