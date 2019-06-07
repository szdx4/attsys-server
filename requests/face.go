package requests

import (
	"encoding/base64"
	"errors"
	"strings"
)

// FaceCreateRequest 更新指定用户人脸信息请求
type FaceCreateRequest struct {
	Info string `binding:"required"`
}

// Validate 验证更新人脸请求的合法性
func (r *FaceCreateRequest) Validate() error {
	tmp := strings.Split(r.Info, "base64,")
	imageBase64 := tmp[1]

	_, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return errors.New("Invalid image format")
	}

	return nil
}
