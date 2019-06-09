package models

// Message 消息模型
type Message struct {
	CommonFields
	Title      string `json:"title"`
	Content    string `json:"content"`
	FromUserID uint   `json:"-"`
	ToUserID   uint   `json:"-"`
	Status     string `json:"status" gorm:"status:enum('unread', 'read')"`
	FromUser   User   `json:"from"`
	ToUser     User   `json:"to"`
}
