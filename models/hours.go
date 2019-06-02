package models

import "time"

// Hours 工时模型
type Hours struct {
	CommonFields
	UserID uint      `json:"user_id"`
	Date   time.Time `json:"date"`
	Hours  int       `json:"hours"`
}
