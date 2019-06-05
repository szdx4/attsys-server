package models

import "time"

// Sign 请假模型
type Sign struct {
	CommonFields
	ShiftID uint      `json:"shift_id"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}
