package models

import (
	"encoding/base64"
	"time"

	"github.com/skip2/go-qrcode"
)

// Qrcode 二维码模型
type Qrcode struct {
	CommonFields
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Image 获取当前二维码的图片
func (m *Qrcode) Image() (string, error) {
	png, err := qrcode.Encode(m.Token, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	str := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)

	return str, nil
}
