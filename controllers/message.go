package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
	"strconv"
)

// MessageShow 获取指定信息
func MessageShow(c *gin.Context) {
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Message ID invalid")
		c.Abort()
		return
	}

	message := models.Message{}
	database.Connector.First(&message, messageID)
	if message.ID < 1 {
		response.NotFound(c, "Message not found")
		c.Abort()
		return
	}

	response.MessageShow(c, message)
}
