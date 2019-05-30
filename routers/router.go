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

	// User
	r.POST("/user/auth", controllers.UserAuth)
	r.GET("/user/:id", middleware.Token, controllers.UserShow)
	r.POST("/user", controllers.UserCreate)
	r.GET("/user", middleware.Token, middleware.Manager, controllers.UserList)
	r.DELETE("/user/:id", middleware.Token, middleware.Manager, controllers.UserDelete)

	// Department
	r.POST("/department", middleware.Token, middleware.Master, controllers.DepartmentCreate)
	r.GET("/department", middleware.Token, controllers.DepartmentList)
	r.PUT("/department/:id", middleware.Token, controllers.DepartmentShow)

	// Face

	// Hours

	// Shift

	// Leave

	// Overtime

	// Sign

	return r
}
