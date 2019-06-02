package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
)

// ShiftCreate 添加排班
func ShiftCreate(c *gin.Context) {
	var req requests.ShiftCreateRequest
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

	shift := models.Shift{
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
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

// UserCreate 新建用户
//func UserCreate(c *gin.Context) {
//
//	user := models.User{
//		Name:         req.Name,
//		Password:     string(hash),
//		DepartmentID: uint(req.Department),
//		Role:         "user",
//	}
//	database.Connector.Create(&user)
//
//	if user.ID < 1 {
//		response.InternalServerError(c, "Internal Server Error")
//		c.Abort()
//		return
//	}
//
//	response.UserCreate(c, user.ID)
//}
