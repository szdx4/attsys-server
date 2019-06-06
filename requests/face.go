package requests

import (
	"encoding/base64"
	"errors"
)

// FaceCreateRequest 更新指定用户人脸信息请求
type FaceCreateRequest struct {
	Info string `binding:"required"`
}

// Validate 验证更新人脸请求的合法性
func (r *FaceCreateRequest) Validate() error {
	_, err := base64.StdEncoding.DecodeString(r.Info)
	if err != nil {
		return errors.New("invalid image format")
	}

	return nil
}
