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
}
