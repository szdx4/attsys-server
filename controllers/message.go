package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/response"
	"github.com/szdx4/attsys-server/utils/database"
)

// MessageShow 获取指定信息
func MessageShow(c *gin.Context) {
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Message ID invalid")
		c.Abort()
		return
	}

	//查找到信息
	message := models.Message{}
	database.Connector.First(&message, messageID)
	if message.ID < 1 {
		response.NotFound(c, "Message not found")
		c.Abort()
		return
	}

	authID, _ := c.Get("user_id")
	if message.FromUserID != authID && message.ToUserID != authID {
		response.Unauthorized(c, "This is not your message")
		c.Abort()
		return
	}

	message.Status = "read"
	database.Connector.Save(&message)

	response.MessageShow(c, message)
}

// MessageList 获取信息列表
func MessageList(c *gin.Context) {
	messages := []models.Message{}
	db := database.Connector
	// 检测 from_user_id
	if fromUserID, isExist := c.GetQuery("from_user_id"); isExist == true {
		fromUserID, _ := strconv.Atoi(fromUserID)
		db = db.Where("from_user_id = ?", fromUserID)
	}

	// 检测 to_user_id
	if toUserID, isExist := c.GetQuery("to_user_id"); isExist == true {
		toUserID, _ := strconv.Atoi(toUserID)
		db = db.Where("to_user_id = ?", toUserID)
	}

	// 检测 status
	if status, isExist := c.GetQuery("status"); isExist == true {
		db = db.Where("status = ?", status)
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
	db.Limit(perPage).Offset((page - 1) * perPage).Find(&messages)
	db.Model(&messages).Count(&total)
	if (page-1)*perPage >= total {
		response.NoContent(c)
		c.Abort()
		return
	}

	response.MessageList(c, total, page, messages)
}
