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

// NoContent 返回 HTTP 204 状态
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest 返回 HTTP 400 状态
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": msg,
	})
}

// Unauthorized 返回 HTTP 401 状态
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  http.StatusUnauthorized,
		"message": msg,
	})
}

// NotFound 返回 HTTP 404 状态
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  http.StatusNotFound,
		"message": msg,
	})
}

// InternalServerError 返回 HTTP 500 状态
func InternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  http.StatusInternalServerError,
		"message": msg,
	})
}
