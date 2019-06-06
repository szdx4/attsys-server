package controllers

import (
	"strconv"

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
