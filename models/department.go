package models

// Department 部门模型
type Department struct {
	CommonFields
	Name      string
	Users     []User `gorm:"foreignkey:DepartmentID"`
	ManagerID uint
}
