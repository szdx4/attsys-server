package models

import (
	"time"
)

// Hours 工时模型
type Hours struct {
	CommonFields
	UserID uint      `json:"-"`
	Date   time.Time `json:"date"`
	Hours  uint      `json:"hours"`
	User   User      `json:"user"`
}
