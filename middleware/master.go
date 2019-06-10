package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/response"
)

// Master 验证用户是否具有 master 的权限
func Master(c *gin.Context) {
	// 获取认证用户角色
	role, ok := c.Get("user_role")
	if !ok {
		response.Unauthorized(c, "Access Denied")
		c.Abort()
		return
	}

	// 判断用户权限
	if role != "master" {
		response.Unauthorized(c, "Access Denied")
		c.Abort()
		return
	}

	c.Next()
}
