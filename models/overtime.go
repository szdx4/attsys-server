package models

import "time"

// Overtime 加班模型
type Overtime struct {
	CommonFields
	UserID  uint      `json:"-"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Remark  string    `json:"remark"`
	Status  string    `json:"status" gorm:"status:enum('wait', 'pass', 'reject')"`
	User    User      `json:"user"`
}
