package models

import "time"

// Sign 请假模型
type Sign struct {
	CommonFields
	ShiftID uint      `json:"-"`        // 请假对应排班
	StartAt time.Time `json:"start_at"` // 请假开始时间
	EndAt   time.Time `json:"end_at"`   // 请假结束时间
	Shift   Shift     `json:"shift"`
}
