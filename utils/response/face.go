package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
)

// FaceShow 获取人脸响应
func FaceShow(c *gin.Context, face models.Face) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   face,
	})
}
