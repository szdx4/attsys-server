package models

import "time"

// shift 排班模型
type Shift struct {
	CommonFields
	UserID  uint      `json:"-"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Type    string    `json:"type" gorm:"type:enum('nomal','overtime','allovertime')"`
	Status  string    `json:"-" gorm:"type:enum('no','on', 'off','leave')"`
}
