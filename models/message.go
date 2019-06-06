package models

// Message 消息模型
type Message struct {
	CommonFields
	Title   string `json:"title"`
	Content string `json:"content"`
	From    uint   `json:"from"`
	To      uint   `json:"to"`
	Status  string `json:"status" gorm:"status:enum('unread', 'read')"`
}
