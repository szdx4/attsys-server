package controllers

import (
	"errors"
	"fmt"
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
	fmt.Println(shift.StartAt)
	fmt.Println(shift.EndAt)
	fmt.Println(shift.Type)

	database.Connector.Create(&shift)
	if shift.ID < 1 {
		response.InternalServerError(c, "Internal Server Error")
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

		if role == "user" && authID != userID {
			response.Unauthorized(c, "You can only get your information")
			c.Abort()
			return
		}

		db = db.Where("user_id = ?", userID)
	} else if role == "user" {
		response.Unauthorized(c, "You can only get your information")
		c.Abort()
		return
	}

	// 检测 start_at
	if startAt, isExit := c.GetQuery("user_id"); isExit {
		db = db.Where("start_at >= ?", startAt)
	}

	// 检测 end_at
	if endAt, isExit := c.GetQuery("user_id"); isExit {
		db = db.Where("end_at <= ?", endAt)
	}

	// 检测 department_id
	if departmentID, isExit := c.GetQuery("department_id"); isExit {
		departmentID, _ := strconv.Atoi(departmentID)

		if role == "manager" {
			manager := models.User{}
			database.Connector.First(&manager, authID)

			if manager.DepartmentID != uint(departmentID) {
				response.Unauthorized(c, "You can only get your department information")
				c.Abort()
				return
			}
		}

		db.Where("users.department_id = ?", departmentID)
	} else if role == "manager" {
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
	var req requests.ShiftDepartmentRequest
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

	// 查找出所有部门员工
	departmentID, _ := strconv.Atoi(c.Param("department_id"))
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
	var shiftIds []uint
	for i := 0; i < len(users); i++ {
		shiftIds = append(shiftIds, users[i].ID)
		shift := models.Shift{
			UserID:  users[i].ID,
			StartAt: startAt,
			EndAt:   endAt,
			Type:    req.Type,
			Status:  "no",
		}
		database.Connector.Create(&shift)
		if shift.ID < 1 {
			response.InternalServerError(c, "Internal Server Error")
			c.Abort()
			return
		}
	}

	response.ShiftDepartment(c, shiftIds)
}

// ShiftUpdate 更新排班状态
func ShiftUpdate(c *gin.Context) {
	var req requests.ShiftUpdateRequest
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
	shiftID, _ := strconv.Atoi(c.Param("shift_id"))
	shift := models.Shift{}
	database.Connector.Where("id = ?", shiftID).First(&shift)
	if shiftID == 0 {
		response.NotFound(c, "shift not found")
		c.Abort()
		return
	}
	// 修改shift的相应信息
	shift.Status = req.Status
	database.Connector.Save(&shift)

	response.ShiftUpdate(c)
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
	database.Connector.Where("id = ?", shiftID).First(&shift)

	if shift.ID == 0 {
		response.NotFound(c, "Shift not found")
		c.Abort()
		return
	}

	database.Connector.Delete(&shift)

	response.ShiftDelete(c)
}
