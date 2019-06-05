package models

import "time"

// Leave 请假模型
type Leave struct {
	CommonFields
	UserID  uint      `json:"user_id"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Remark  string    `json:"remark"`
	Status  string    `json:"status" gorm:"status:enum('wait', 'pass', 'reject')"`
}
