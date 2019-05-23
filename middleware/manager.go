package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/utils/response"
)

// Manager 验证用户是否具有 manager 以上的权限
func Manager(c *gin.Context) {
	role, ok := c.Get("user_role")
	if !ok {
		response.Unauthorized(c, "Access Denied")
		c.Abort()
		return
	}

	if role != "manager" && role != "master" {
		response.Unauthorized(c, "Access Denied")
		c.Abort()
		return
	}

	c.Next()
}
