package models

import "time"

// Shift 排班模型
type Shift struct {
	CommonFields
	UserID  uint      `json:"-"`                                                   // 排班所属用户 ID
	StartAt time.Time `json:"start_at"`                                            // 排班开始时间
	EndAt   time.Time `json:"end_at"`                                              // 排班结束时间
	Type    string    `json:"type" gorm:"type:enum('normal', 'allovertime')"`      // 排班类型
	Status  string    `json:"status" gorm:"type:enum('no', 'on', 'off', 'leave')"` // 排班审核状态
	User    User      `json:"user"`
}
