package message

import (
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
)

// Send 发送信息
func Send(fromID, toID uint, title, content string) {
	message := models.Message{
		Title:      title,
		Content:    content,
		FromUserID: fromID,
		ToUserID:   toID,
		Status:     "unread",
	}
	database.Connector.Create(&message)
}
