package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/routers"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/qcloud"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 初始化数据库和数据表
	database.Connect()
	database.Migrate()
	database.Seed()

	// 初始化腾讯云人脸识别库
	qcloud.GroupInit()
	qcloud.PersonInit()

	// 初始化路由
	router := routers.Router()

	// 初始化 HTTP 服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Server.HTTPPort),
		Handler:        router,
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动 HTTP 服务器
	log.Printf("Server listening at port: %d", config.Server.HTTPPort)
	server.ListenAndServe()
}
