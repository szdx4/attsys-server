package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/utils/response"
)

// Master 验证用户是否具有 master 的权限
func Master(c *gin.Context) {
	role, ok := c.Get("user_role")
	if !ok {
		response.Unauthorized(c, "Access Denied")
		c.Abort()
		return
	}

	if role != "master" {
		response.Unauthorized(c, "Access Denied")
		c.Abort()
		return
	}

	c.Next()
}
