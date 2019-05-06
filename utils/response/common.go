package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK 返回 HTTP 200 状态
func OK(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": msg,
		"data":    data,
	})
}

// Created 返回 HTTP 201 状态
func Created(c *gin.Context, id uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"resourceId": id,
	})
}

// BadRequest 返回 HTTP 400 状态
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": msg,
	})
}
