package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/qcloud"
)

// FaceUserShow 获取指定用户可用的人脸信息
func FaceUserShow(c *gin.Context) {
	// 从 URL 中获取用户 ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	// 验证权限
	authID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	if role != "master" && authID != userID {
		response.Unauthorized(c, "Unauthorized")
		c.Abort()
		return
	}

	// 从数据库中查询人脸信息
	face := models.Face{}
	database.Connector.Preload("User").Where("user_id = ? AND status = 'available'", userID).First(&face)
	if face.ID == 0 {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
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

	// 验证提交数据的合法性
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err.Error())
		c.Abort()
		return
	}

	// 从 URL 中获取用户 ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "User ID invalid")
		c.Abort()
		return
	}

	// 判断人脸操作权限
	authID, _ := c.Get("user_id")
	if userID != authID.(int) {
		response.Unauthorized(c, "You can only update face info for yourself")
		c.Abort()
		return
	}

	// 更新人脸信息
	face := models.Face{
		UserID: uint(userID),
		Info:   req.Info,
		Status: "wait",
	}
	database.Connector.Create(&face)

	// 检验是否操作成功
	if face.ID == 0 {
		response.InternalServerError(c, "Database error")
		c.Abort()
		return
	}

	// 发送响应
	response.FaceCreate(c, face.ID)
}

// FaceList 获取人脸列表
func FaceList(c *gin.Context) {
	faces := []models.Face{}
	db := database.Connector.Preload("User")

	// 按照用户 ID 过滤
	if userID, isExit := c.GetQuery("user_id"); isExit {
		userID, _ := strconv.Atoi(userID)
		db = db.Where("user_id = ?", userID)
	}

	// 按照人脸状态过滤
	if status, isExit := c.GetQuery("status"); isExit {
		status, _ := strconv.Atoi(status)
		db = db.Where("status = ?", status)
	}

	// 处理分页
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	perPage := config.App.ItemsPerPage
	total := 0

	// 执行查询
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&faces)
	db.Model(&models.Face{}).Count(&total)

	// 判断当前页是否为空
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	// 发送响应
	response.FaceList(c, total, page, perPage, faces)
}

// FaceUpdate 编辑人脸信息
func FaceUpdate(c *gin.Context) {
	// 从 URL 中获取人脸 ID
	faceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Face ID invalid")
		c.Abort()
		return
	}

	// 从数据库中查询人脸信息
	face := models.Face{}
	database.Connector.First(&face, faceID)
	if face.ID == 0 {
		response.NotFound(c, "Face not found")
		c.Abort()
		return
	}

	// 判断人脸状态
	if face.Status != "wait" {
		response.BadRequest(c, "Face status invalid")
		c.Abort()
		return
	}

	// 禁用该用户其他可用人脸信息
	faces := []models.Face{}
	database.Connector.Where("user_id = ? AND status = 'available'", face.UserID).Find(&faces)
	for _, item := range faces {
		item.Status = "discarded"
		database.Connector.Save(&item)
	}

	// 将此人脸信息设置为可用
	face.Status = "available"
	database.Connector.Save(&face)

	// 更新腾讯云人脸库中的信息
	qcloud.DeletePersonFromGroup(config.Qcloud.GroupName, strconv.Itoa(int(face.UserID)))
	qcloud.CreatePerson(config.Qcloud.GroupName, strconv.Itoa(int(face.UserID)), face.Info)

	// 发送响应
	response.FaceUpdate(c)
}
