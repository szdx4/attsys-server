package main

import (
	"fmt"
	"net/http"

	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/routers"
)

func main() {
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
