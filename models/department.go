package models

// Department 部门模型
type Department struct {
	CommonFields
	Name      string `json:"name"`
	Users     []User `gorm:"foreignkey:DepartmentID"`
	ManagerID uint   `json:"manager"`
}
