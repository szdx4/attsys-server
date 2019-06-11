package models

// User 用户模型
type User struct {
	CommonFields
	Name         string     `json:"name"`                                             // 用户名
	Password     string     `json:"-"`                                                // 用户密码
	Role         string     `json:"role" gorm:"type:enum('user','manager','master')"` // 用户角色
	DepartmentID uint       `json:"-"`                                                // 用户所属部门 ID
	Hours        uint       `json:"hours"`                                            // 用户总工时
	Department   Department `json:"department"`                                       // 用户所属部门
}
