package routers

import (
	"github.com/szdx4/attsys-server/controllers"
	"github.com/szdx4/attsys-server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
)

// Router 设置路由和公共中间件，返回 Gin Engine 对象
func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(config.App.RunMode)

	r.GET("/", controllers.Home)
	r.POST("/user/auth", controllers.UserAuth)
	r.GET("/user/:id", middleware.Token, controllers.UserShow)
	r.POST("/user", controllers.UserCreate)

	return r
}
