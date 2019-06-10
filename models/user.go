package models

// User 用户模型
type User struct {
	CommonFields
	Name         string     `json:"name"`                                             // 用户名
	Password     string     `json:"-"`                                                // 用户密码
	Role         string     `json:"role" gorm:"type:enum('user','manager','master')"` // 用户角色
	DepartmentID uint       `json:"-"`                                                // 用户所属部门
	Hours        uint       `json:"hours"`                                            // 用户总工时
	Department   Department `json:"department"`
	Shifts       []Shift    `json:"-" gorm:"foreignkey:UserID"`
	Hourss       []Hours    `json:"-" gorm:"foreignkey:UserID"`
	Leaves       []Leave    `json:"-" gorm:"foreignkey:UserID"`
	Overtimes    []Overtime `json:"-" gorm:"foreignkey:UserID"`
	MessageFrom  []Message  `json:"-" gorm:"foreignkey:FromUserID"`
	MessageTo    []Message  `json:"-" gorm:"foreignkey:ToUserID"`
}
