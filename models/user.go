package models

import (
	"github.com/jinzhu/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Name       string
	Password   string
	Role       string
	Department uint
	Hours      uint
}
