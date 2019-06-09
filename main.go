package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/szdx4/attsys-server/utils/qcloud"

	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/routers"
	"github.com/szdx4/attsys-server/utils/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database.Connect()
	database.Migrate()
	database.Seed()

	qcloud.GroupInit()
	qcloud.PersonInit()

	router := routers.Router()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Server.HTTPPort),
		Handler:        router,
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Server listening at port: %d", config.Server.HTTPPort)
	server.ListenAndServe()
}
