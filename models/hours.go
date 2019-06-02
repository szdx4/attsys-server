package models

import "time"

// Hours 工时模型
type Hours struct {
	CommonFields
	User  uint
	Date  time.Time `json:"date"`
	Hours int
}
