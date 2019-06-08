package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
)

// ShiftCreate 添加排班
func ShiftCreate(c *gin.Context) {
	var req requests.ShiftCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	if err := req.Validate(c); err != nil {
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
	shift := models.Shift{
		UserID:  uint(userID),
		StartAt: startAt,
		EndAt:   endAt,
		Type:    req.Type,
		Status:  "no",
	}

	database.Connector.Create(&shift)
	if shift.ID < 1 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	response.ShiftCreate(c, shift.ID)
}

// ShiftList 排班列表
func ShiftList(c *gin.Context) {
	shifts := []models.Shift{}
	db := database.Connector.Joins("LEFT JOIN users ON shifts.user_id = users.id")

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	// 检测 user_id
	if userID, isExit := c.GetQuery("user_id"); isExit {
		userID, _ := strconv.Atoi(userID)

		// user 只能查询自己的排班
		if role == "user" && authID != userID {
			response.Unauthorized(c, "You can only get your information")
			c.Abort()
			return
		}

		db = db.Where("user_id = ?", userID)
	} else if role == "user" {
		// user 不能得到列表
		response.Unauthorized(c, "You can only get your information")
		c.Abort()
		return
	}

	// 检测 start_at
	if startAt, isExit := c.GetQuery("start_at"); isExit {
		db = db.Where("start_at >= ?", startAt)
	}

	// 检测 end_at
	if endAt, isExit := c.GetQuery("end_at"); isExit {
		db = db.Where("end_at <= ?", endAt)
	}

	// 检测 department_id
	if departmentID, isExit := c.GetQuery("department_id"); isExit {
		departmentID, _ := strconv.Atoi(departmentID)

		// manager 只能查看自己的部门
		if role == "manager" {
			manager := models.User{}
			database.Connector.First(&manager, authID)

			if manager.DepartmentID != uint(departmentID) {
				response.Unauthorized(c, "You can only get your department information")
				c.Abort()
				return
			}
		}

		db = db.Where("users.department_id = ?", departmentID)
	} else if role == "manager" {
		// manager不能得到全体列表
		response.Unauthorized(c, "You can only get your department information")
		c.Abort()
		return
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

	db.Limit(perPage).Offset((page - 1) * perPage).Find(&shifts)
	db.Model(&models.Shift{}).Count(&total)

	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.ShiftList(c, total, page, shifts)
}

// ShiftDepartment 部门排班
func ShiftDepartment(c *gin.Context) {
	departmentID, _ := strconv.Atoi(c.Param("department_id"))

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	var req requests.ShiftDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	if err := req.Validate(departmentID, role.(string), authID.(int)); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 查找出所有部门员工
	users := []models.User{}
	database.Connector.Where("department_id = ?", departmentID).Find(&users)

	// 获得起始时间
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

	// 存入数据库
	var shiftIDs []uint
	for _, user := range users {
		shift := models.Shift{
			UserID:  user.ID,
			StartAt: startAt,
			EndAt:   endAt,
			Type:    req.Type,
			Status:  "no",
		}
		database.Connector.Create(&shift)
		if shift.ID == 0 {
			response.InternalServerError(c, "Database error")
			c.Abort()
			return
		}
		shiftIDs = append(shiftIDs, user.ID)
	}

	response.ShiftDepartment(c, shiftIDs)
}

// ShiftDelete 删除排班
func ShiftDelete(c *gin.Context) {
	shiftID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Shift ID invalid")
		c.Abort()
		return
	}
	shift := models.Shift{}
	database.Connector.First(&shift, shiftID).Related(&shift.User)

	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	role, _ := c.Get("user_role")
	authID, _ := c.Get("user_id")

	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		if manager.DepartmentID != shift.User.DepartmentID {
			response.Unauthorized(c, "You can only delete your department shifts")
			c.Abort()
			return
		}
	}

	database.Connector.Delete(&shift)

	response.ShiftDelete(c)
}
