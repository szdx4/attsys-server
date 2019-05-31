package models

// Hours 工时模型
type Hours struct {
	CommonFields
	User  uint
	Date  string
	Hours int
}
