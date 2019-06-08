package controllers

import (
	"strconv"
	"time"

	"github.com/szdx4/attsys-server/utils/message"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
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

	authID, _ := c.Get("user_id")
	userID, _ := strconv.Atoi(c.Param("id"))

	if authID != userID {
		response.BadRequest(c, "You can only apply overtime for yourself")
		c.Abort()
		return
	}

	user := models.User{}
	database.Connector.First(&user, userID)

	shift := models.Shift{}
	database.Connector.Where("end_at < ? AND status = 'off'", time.Now()).Order("end_at DESC").First(&shift)
	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	overtime := models.Overtime{
		UserID:  uint(userID),
		StartAt: shift.EndAt,
		EndAt:   time.Now(),
		Remark:  req.Remark,
		Status:  "wait",
	}
	database.Connector.Create(&overtime)
	if overtime.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	message.Send(user.ID, user.Department.ManagerID, "加班申请", "理由："+overtime.Remark)

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

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 用户只能获取自己的加班
	if role == "user" && userID != authID {
		response.Unauthorized(c, "You can only get your own overtime")
		c.Abort()
		return
	}

	// 部门主管只能获得本部门的员工加班
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		user := models.User{}
		database.Connector.First(&user, userID)
		if manager.DepartmentID != user.DepartmentID {
			response.Unauthorized(c, "You can only get your department overtime")
			c.Abort()
			return
		}
	}

	// 查找加班
	overtime := []models.Overtime{}
	db := database.Connector.Where("user_id = ?", userID)
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&overtime)
	db.Model(&models.Overtime{}).Count(&total)

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
	db := database.Connector.Joins("LEFT JOIN users ON users.id = overtimes.user_id").Order("created_at DESC")

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 部门主管只能获得自己部门列表
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		db = db.Where("users.department_id = ?", manager.DepartmentID)
	}

	db.Limit(perPage).Offset((page - 1) * perPage).Find(&overtime)
	db.Model(&models.Overtime{}).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.OvertimeList(c, total, page, overtime)
}

// OvertimeUpdate 审批加班
func OvertimeUpdate(c *gin.Context) {
	var req requests.OvertimeUpdateRequest
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

	// 查找指定加班
	overtimeID, _ := strconv.Atoi(c.Param("id"))
	overtime := models.Overtime{}
	database.Connector.Preload("User").First(&overtime, overtimeID)
	if overtime.ID == 0 {
		response.NotFound(c, "Overtime not found")
		c.Abort()
		return
	}

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 主管只能给本部门加班审核
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		if manager.DepartmentID != overtime.User.DepartmentID {
			response.Unauthorized(c, "You can only edit your department overtime")
			c.Abort()
			return
		}
	}

	// 修改 overtime 的 status
	overtime.Status = req.Status
	database.Connector.Save(&overtime)

	if overtime.Status == "pass" {
		// 创建工时记录
		hour := uint(overtime.EndAt.Sub(overtime.StartAt).Hours())
		hours := models.Hours{
			UserID: overtime.UserID,
			Date:   overtime.EndAt,
			Hours:  hour,
		}
		database.Connector.Create(&hours)
		hours.User.Hours += hour
		database.Connector.Save(&hours.User)

		message.Send(uint(authID.(int)), overtime.UserID, "加班申请结果", "加班申请通过")
	} else {
		message.Send(uint(authID.(int)), overtime.UserID, "加班申请结果", "加班申请未通过")
	}

	response.OvertimeUpdate(c)
}
