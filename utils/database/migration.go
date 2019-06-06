package database

import "github.com/szdx4/attsys-server/models"

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
