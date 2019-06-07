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

// FaceCreate 更新人脸信息响应
func FaceCreate(c *gin.Context, faceID uint) {
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"face_id": faceID,
	})
}

// FaceList 人脸列表响应
func FaceList(c *gin.Context, total, page, perPage int, faces []models.Face) {
	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total":        total,
		"current_page": page,
		"per_page":     perPage,
		"data":         faces,
	})
}

// FaceUpdate 编辑人脸响应
func FaceUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
