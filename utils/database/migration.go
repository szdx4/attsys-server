package database

import "github.com/szdx4/attsys-server/models"

// Migrate 执行数据库迁移
func Migrate() {
	Connector.AutoMigrate(&models.Department{})
	Connector.AutoMigrate(&models.User{})
}
