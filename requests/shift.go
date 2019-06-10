package requests

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/common"
	"github.com/szdx4/attsys-server/utils/database"
)

// ShiftCreateRequest 添加排班
type ShiftCreateRequest struct {
	StartAt string `binding:"required" json:"start_at"`
	EndAt   string `binding:"required" json:"end_at"`
	Type    string `binding:"required"`
}

// Validate 验证 ShiftCreateRequest 请求中的有效性
func (r *ShiftCreateRequest) Validate(c *gin.Context) error {
	// 将接收的 string 格式转换成 Time
	startAt, err := common.ParseTime(r.StartAt)
	if err != nil {
		return errors.New("start_at not valid")
	}
	endAt, err := common.ParseTime(r.EndAt)
	if err != nil {
		return errors.New("end_at not valid")
	}

	// 验证给出排班的先后有效性
	if startAt.After(endAt) {
		return errors.New("Time not valid")
	}

	// 验证个人排班的冲突性
	userID, _ := strconv.Atoi(c.Param("id"))

	user := models.User{}
	database.Connector.First(&user, userID)
	if user.ID == 0 {
		return errors.New("User not found")
	}

	shift := models.Shift{}
	db := database.Connector
	db = db.Where("user_id = ?", userID)
	db = db.Where("start_at < ?", endAt)
	db = db.Where("end_at > ?", startAt)
	db.First(&shift)
	if shift.ID > 0 {
		return errors.New("Time is conflicting")
	}

	// 验证排班类型的有效性
	if r.Type != "normal" && r.Type != "allovertime" {
		return errors.New("Type not valid")
	}

	// 无误则返回空
	return nil
}

// ShiftDepartmentRequest 部门排班
type ShiftDepartmentRequest struct {
	StartAt string `binding:"required" json:"start_at"`
	EndAt   string `binding:"required" json:"end_at"`
	Type    string `binding:"required"`
}

// Validate 验证 ShiftDepartmentRequest 请求中的有效性
func (r *ShiftDepartmentRequest) Validate(departmentID int, role string, authID int) error {
	// 将接收的 string 格式转换成 Time
	startAt, err := common.ParseTime(r.StartAt)
	if err != nil {
		return errors.New("start_at not valid")
	}
	endAt, err := common.ParseTime(r.EndAt)
	if err != nil {
		return errors.New("end_at not valid")
	}

	// 验证给出排班的有效性
	if startAt.After(endAt) {
		return errors.New("Time not valid")
	}

	// 验证类型的有效性
	if r.Type != "normal" && r.Type != "allovertime" {
		return errors.New("Type not valid")
	}

	department := models.Department{}
	database.Connector.First(&department, departmentID)
	if department.ID == 0 {
		return errors.New("Department not found")
	}

	// 验证
	if role == "manager" {
		manager := models.User{}
		database.Connector.First(&manager, authID)
		if manager.DepartmentID != uint(departmentID) {
			return errors.New("You can only arrange your department shifts")
		}
	}

	return nil
}
