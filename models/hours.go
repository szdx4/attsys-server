package models

import "time"

// Hours 工时模型
type Hours struct {
	CommonFields
	UserID uint      `json:"user_id"`
	Date   time.Time `json:"date"`
	Hours  uint      `json:"hours"`
}

// Hours 响应结构
type HourData struct {
	ID       uint      `json:"id"`
	UserID   uint      `json:"user_id"`
	UserName string    `json:"user_name"`
	Date     time.Time `json:"date"`
	Hours    int       `json:"hours"`
}
