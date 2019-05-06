package models

import (
	"github.com/jinzhu/gorm"
)

// Department 部门模型
type Department struct {
	gorm.Model
	Name  string
	Users []User `gorm:"foreignkey:Department"`
}
