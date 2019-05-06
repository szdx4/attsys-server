package main

import (
	"fmt"
	"net/http"

	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/routers"
	"github.com/szdx4/attsys-server/utils/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database.Connect()
	database.Migrate()

	router := routers.Router()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Server.HTTPPort),
		Handler:        router,
		ReadTimeout:    config.Server.ReadTimeout,
		WriteTimeout:   config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
