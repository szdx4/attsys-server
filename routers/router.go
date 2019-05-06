package routers

import (
	"net/http"

	"github.com/szdx4/attsys-server/utils/database"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/requests"
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

	r.POST("/department", func(c *gin.Context) {
		var req requests.CreateDepartmentRequest
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "error",
			})
			return
		}

		department := models.Department{
			Name: req.Name,
		}
		database.Connector.Create(&department)

		c.JSON(http.StatusCreated, gin.H{
			"status":     http.StatusCreated,
			"resourseId": department.ID,
		})
	})

	return r
}
