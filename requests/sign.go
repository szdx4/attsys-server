package requests

import (
	"errors"
	"time"

	"github.com/szdx4/attsys-server/models"
	"github.com/szdx4/attsys-server/utils/database"
)

// SignWithQrcodeRequest 二维码签到请求
type SignWithQrcodeRequest struct {
	Token string `json:"token"`
}

// Validate 验证二维码签到请求的有效性
func (r *SignWithQrcodeRequest) Validate() error {
	qrcode := models.Qrcode{}
	database.Connector.Where("token = ?", r.Token).First(&qrcode)

	if qrcode.ID == 0 {
		return errors.New("Qrcode not found")
	}

	if time.Now().After(qrcode.ExpiredAt) {
		return errors.New("Qrcode expired")
	}

	return nil
}
