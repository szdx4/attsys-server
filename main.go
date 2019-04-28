package main

import (
	"fmt"
	"net/http"

	"github.com/szdx4/attsys-server/routers"
	"github.com/szdx4/attsys-server/utils/setting"
)

func main() {
	router := routers.Router()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
