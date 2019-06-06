package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
	"github.com/szdx4/attsys-server/utils/response"
)

// SignQrcode 获取二维码
func SignQrcode(c *gin.Context) {
	qrcode := models.Qrcode{
		ExpiredAt: time.Now().Add(60 * time.Second),
	}
	qrcode.RandToken()
	database.Connector.Create(&qrcode)

	if qrcode.ID == 0 {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	image, err := qrcode.Image()
	if err != nil {
		response.InternalServerError(c, "Internal Server Error")
		c.Abort()
		return
	}

	response.SignQrcode(c, image, qrcode.ExpiredAt)
}
