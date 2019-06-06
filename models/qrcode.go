package models

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/skip2/go-qrcode"
)

// Qrcode 二维码模型
type Qrcode struct {
	CommonFields
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

// RandToken 随机生成 Token
func (m *Qrcode) RandToken() {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 128)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	m.Token = string(b)
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
