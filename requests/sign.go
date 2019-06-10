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
	// 验证二维码是否存在
	qrcode := models.Qrcode{}
	database.Connector.Where("token = ?", r.Token).First(&qrcode)
	if qrcode.ID == 0 {
		return errors.New("Qrcode not found")
	}

	// 验证二维码是否过期
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
	// 验证人脸是否匹配
	resID, err := qcloud.SearchFaces(config.Qcloud.GroupName, r.Face)
	if err != nil {
		return errors.New("Face not match")
	}
	if resID != strconv.Itoa(userID) {
		return errors.New("Face not match")
	}

	return nil
}
