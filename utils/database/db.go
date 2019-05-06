package database

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/szdx4/attsys-server/config"
)

// Connector 数据库连接器
var Connector *gorm.DB

// Connect 连接数据库
func Connect() {
	var err error
	Connector, err = gorm.Open(config.Database.Type, config.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
}
