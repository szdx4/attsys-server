package routers

import (
	"net/http"

	"github.com/szdx4/attsys-server/utils/setting"

	"github.com/gin-gonic/gin"
)

// Router 设置路由和公共中间件，返回 Gin Engine 对象
func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "hello world",
		})
	})

	return r
}
