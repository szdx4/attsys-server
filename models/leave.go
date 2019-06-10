package models

import "time"

// Leave 请假模型
type Leave struct {
	CommonFields
	UserID  uint      `json:"-"`                                                                // 请假申请用户
	StartAt time.Time `json:"start_at"`                                                         // 请假开始时间
	EndAt   time.Time `json:"end_at"`                                                           // 请假结束时间
	Remark  string    `json:"remark"`                                                           // 请假原因
	Status  string    `json:"status" gorm:"status:enum('wait', 'pass', 'reject', 'discarded')"` // 请假审核状态
	User    User      `json:"user"`
}
