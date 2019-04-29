package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
)

// Router 设置路由和公共中间件，返回 Gin Engine 对象
func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(config.App.RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "hello world",
		})
	})

	return r
}
