package models

import "time"

// Overtime 加班模型
type Overtime struct {
	CommonFields
	UserID  uint      `json:"-"`                                                   // 加班申请用户 ID
	StartAt time.Time `json:"start_at"`                                            // 加班开始时间
	EndAt   time.Time `json:"end_at"`                                              // 加班结束时间
	Remark  string    `json:"remark"`                                              // 加班原因
	Status  string    `json:"status" gorm:"status:enum('wait', 'pass', 'reject')"` // 加班审核状态
	User    User      `json:"user"`
}
