package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home 默认页面
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "attsys-server",
	})
}
