package models

import "time"

// Shift 排班模型
type Shift struct {
	CommonFields
	UserID  uint      `json:"user_id"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Type    string    `json:"type" gorm:"type:enum('normal','overtime','allovertime')"`
	Status  string    `json:"status" gorm:"type:enum('no','on', 'off','leave')"`
	User    User      `json:"-"`
}
