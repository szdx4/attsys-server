package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
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
	if role != "master" && authID != userID {
		response.Unauthorized(c, "Unauthorized")
		c.Abort()
		return
	}

	face := models.Face{}
	database.Connector.Preload("User").Where("user_id = ? AND status = 'available'", userID).First(&face)

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

	authID, _ := c.Get("user_id")

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	if userID != authID.(int) {
		response.Unauthorized(c, "You can only update face info for yourself")
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
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	response.FaceCreate(c, face.ID)
}

// FaceList 获取人脸列表
func FaceList(c *gin.Context) {
	faces := []models.Face{}
	db := database.Connector.Preload("User")

	if userID, isExit := c.GetQuery("user_id"); isExit {
		userID, _ := strconv.Atoi(userID)
		db = db.Where("user_id = ?", userID)
	}

	if status, isExit := c.GetQuery("status"); isExit {
		status, _ := strconv.Atoi(status)
		db = db.Where("status = ?", status)
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	total := 0

	db.Limit(perPage).Offset((page - 1) * perPage).Find(&faces)
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}
	db.Model(&models.Face{}).Count(&total)

	response.FaceList(c, total, page, perPage, faces)
}

// FaceUpdate 编辑人脸信息
func FaceUpdate(c *gin.Context) {
	faceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Face ID invalid")
		c.Abort()
		return
	}

	face := models.Face{}
	database.Connector.First(&face, faceID)
	if face.ID == 0 {
		response.NotFound(c, "Face not found")
		c.Abort()
		return
	}

	if face.Status != "wait" {
		response.BadRequest(c, "Face status invalid")
		c.Abort()
		return
	}

	faces := []models.Face{}
	database.Connector.Where("user_id = ? AND status = 'available'", face.UserID).Find(&faces)

	for _, item := range faces {
		item.Status = "discarded"
		database.Connector.Save(&item)
	}

	face.Status = "available"
	database.Connector.Save(&face)

	response.FaceUpdate(c)
}
