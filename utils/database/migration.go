package database

import (
	"github.com/szdx4/attsys-server/models"
	"golang.org/x/crypto/bcrypt"
)

// Migrate 执行数据库迁移
func Migrate() {
	Connector.AutoMigrate(&models.Department{})
	Connector.AutoMigrate(&models.User{})
	Connector.AutoMigrate(&models.Face{})
	Connector.AutoMigrate(&models.Shift{})
	Connector.AutoMigrate(&models.Hours{})
	Connector.AutoMigrate(&models.Leave{})
	Connector.AutoMigrate(&models.Overtime{})
	Connector.AutoMigrate(&models.Sign{})
	Connector.AutoMigrate(&models.Qrcode{})
	Connector.AutoMigrate(&models.Message{})
}

// Seed 执行数据库填充
func Seed() {
	userCount := 0
	Connector.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("root"), 10)
		if err != nil {
			panic(err)
		}

		user := models.User{}
		user.Name = "root"
		user.Password = string(hash)
		user.DepartmentID = 0
		user.Role = "master"

		Connector.Save(&user)
	}
}
