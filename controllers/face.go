package controllers

import (
	"strconv"

	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/response"
)

// FaceUserShow 获取指定用户可用的人脸信息
func FaceUserShow(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	authID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	if role == "user" && authID != userID {
		response.Unauthorized(c, "Unauthorized")
		c.Abort()
		return
	}

	face := models.Face{}
	database.Connector.Where("user_id = ? AND status = 'available'", userID).First(&face)

	if face.ID == 0 {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.FaceShow(c, face)
}

// FaceCreate 更新指定用户人脸信息
func FaceCreate(c *gin.Context) {
	var req requests.FaceCreateRequest
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

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	face := models.Face{
		UserID: uint(userID),
		Info:   req.Info,
		Status: "wait",
	}
	database.Connector.Create(&face)

	if face.ID == 0 {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	response.FaceCreate(c, face.ID)
}
