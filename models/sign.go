package models

import "time"

// Sign 签到模型
type Sign struct {
	CommonFields
	ShiftID uint      `json:"-"`        // 签到对应排班
	StartAt time.Time `json:"start_at"` // 签到开始时间
	EndAt   time.Time `json:"end_at"`   // 签到结束时间
	Shift   Shift     `json:"shift"`
}
