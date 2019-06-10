package models

// Message 消息模型
type Message struct {
	CommonFields
	Title      string `json:"title"`                                       // 消息标题
	Content    string `json:"content"`                                     // 消息内容
	Status     string `json:"status" gorm:"status:enum('unread', 'read')"` // 消息读取状态
	FromUserID uint   `json:"-"`                                           // 消息发送人
	ToUserID   uint   `json:"-"`                                           // 消息接收人
	FromUser   User   `json:"from"`
	ToUser     User   `json:"to"`
}
