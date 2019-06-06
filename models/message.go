package models

// Message 消息模型
type Message struct {
	CommonFields
	Title      string `json:"title"`
	Content    string `json:"content"`
	FromUserID uint   `json:"from_user_id"`
	ToUserID   uint   `json:"to_user_id"`
	Status     string `json:"status" gorm:"status:enum('unread', 'read')"`
}
