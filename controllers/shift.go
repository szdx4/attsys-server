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
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	response.ShiftCreate(c, shift.ID)
}

// ShiftList 排班列表
func ShiftList(c *gin.Context) {
	shifts := []models.Shift{}
	db := database.Connector
	// 检测 user_id
	if userID, isExit := c.GetQuery("user_id"); isExit == true {
		userID, _ := strconv.Atoi(userID)
		db = db.Where("user_id = ?", userID)
	}

	// 检测 department_id
	//if departmentID, isExit := c.GetQuery("department_id"); isExit == true {
	//	departmentID, _ := strconv.Atoi(departmentID)
	//	db = db.Preload("User", "department_id = ?", departmentID)
	//}

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

	// 用遍历的方法检测 department_id
	if departmentID, isExit := c.GetQuery("department_id"); isExit == true {
		departmentID, _ := strconv.Atoi(departmentID)
		for i := 0; i < len(shifts); {
			user := models.User{}
			check := database.Connector
			check.Where("id = ?", shifts[i].UserID).First(&user)
			if user.DepartmentID != uint(departmentID) {
				shifts = append(shifts[:i], shifts[i+1:]...)
			} else {
				i++
			}
		}
	}

	db.Model(&shifts).Count(&total)
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

	if err := req.Validate(c); err != nil {
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
