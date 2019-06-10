package models

// Department 部门模型
type Department struct {
	CommonFields
	Name string `json:"name"` // 部门名称
}
