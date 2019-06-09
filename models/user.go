package models

// User 用户模型
type User struct {
	CommonFields
	Name         string     `json:"name"`
	Password     string     `json:"-"`
	Role         string     `json:"role" gorm:"type:enum('user','manager','master')"`
	DepartmentID uint       `json:"-"`
	Department   Department `json:"department"`
	Hours        uint       `json:"hours"`
	Shifts       []Shift    `json:"-" gorm:"foreignkey:UserID"`
	Hourss       []Hours    `json:"-" gorm:"foreignkey:UserID"`
	Leaves       []Leave    `json:"-" gorm:"foreignkey:UserID"`
	Overtimes    []Overtime `json:"-" gorm:"foreignkey:UserID"`
	MessageFrom  []Message  `json:"-" gorm:"foreignkey:FromUserID"`
	MessageTo    []Message  `json:"-" gorm:"foreignkey:ToUserID"`
}
