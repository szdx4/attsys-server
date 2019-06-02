package models

// Department 部门模型
type Department struct {
	CommonFields
	Name      string `json:"name"`
	ManagerID uint   `json:"manager"`
	Users     []User `json:"-" gorm:"foreignkey:DepartmentID"`
}
