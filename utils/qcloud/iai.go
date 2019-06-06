package qcloud

import (
	"strings"

	"github.com/szdx4/attsys-server/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	iai "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/iai/v20180301"
)

// CompareFace 人脸比对
func CompareFace(faceA, faceB string) (float64, error) {
	if strings.Index(faceA, "data:") > -1 {
		tmp := strings.Split(faceA, "base64,")
		faceA = tmp[1]
	}
	if strings.Index(faceB, "data:") > -1 {
		tmp := strings.Split(faceB, "base64,")
		faceB = tmp[1]
	}

	credential := common.NewCredential(
		config.Qcloud.SecretID,
		config.Qcloud.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 10
	cpf.HttpProfile.Endpoint = "iai.tencentcloudapi.com"
	cpf.SignMethod = "TC3HmacSHA256"

	client, _ := iai.NewClient(credential, regions.Beijing, cpf)

	request := iai.NewCompareFaceRequest()
	request.ImageA = common.StringPtr(faceA)
	request.ImageB = common.StringPtr(faceB)

	response, err := client.CompareFace(request)
	if err != nil {
		return 0.0, err
	}

	return *response.Response.Score, nil
}
