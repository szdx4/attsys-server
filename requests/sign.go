package requests

import (
	"errors"
	"strconv"
	"time"

	"github.com/szdx4/attsys-server/config"
	"github.com/szdx4/attsys-server/utils/qcloud"

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

// SignWithFaceRequest 人脸签到请求
type SignWithFaceRequest struct {
	Face string `json:"face"`
}

// Validate 验证人脸签到请求的合法性
func (r *SignWithFaceRequest) Validate(userID int) error {
	face := models.Face{}
	database.Connector.Where("user_id = ? AND status = 'available'", userID).First(&face)

	if face.ID == 0 {
		return errors.New("Face info not found")
	}

	resID, err := qcloud.SearchFaces(config.Qcloud.GroupName, r.Face)
	if err != nil {
		return errors.New("Face not match")
	}
	if resID != strconv.Itoa(userID) {
		return errors.New("Face not match")
	}

	return nil
}
