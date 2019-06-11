package models

import (
	"time"
)

// Hours 工时模型
type Hours struct {
	CommonFields
	UserID uint      `json:"-"`     // 工时所属用户
	Date   time.Time `json:"date"`  // 工时添加日期
	Hours  uint      `json:"hours"` // 工时数量（h）
	User   User      `json:"user"`
}
