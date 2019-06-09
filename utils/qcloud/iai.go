package qcloud

import (
	goerrors "errors"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	iai "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/iai/v20180301"
)

func getQcloudIaiClient() *iai.Client {
	credential := common.NewCredential(
		"AKIDOaaKvJVXM8g28Bycx9KI3ZsBIAsQlrjO",
		"ax4wjJisrcATWVGo7DF2at43oSu8AwUv",
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 10
	cpf.HttpProfile.Endpoint = "iai.tencentcloudapi.com"
	cpf.SignMethod = "TC3HmacSHA256"

	client, _ := iai.NewClient(credential, regions.Beijing, cpf)

	return client
}

func getImageBase64(image string) string {
	if strings.Index(image, "base64,") > -1 {
		tmp := strings.Split(image, "base64,")
		return tmp[1]
	}

	return image
}

// DeleteGroup 删除人员库
func DeleteGroup(groupID string) error {
	client := getQcloudIaiClient()

	request := iai.NewDeleteGroupRequest()
	request.GroupId = common.StringPtr(groupID)

	_, err := client.DeleteGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

// CreateGroup 创建人员库
func CreateGroup(groupID string) error {
	client := getQcloudIaiClient()

	request := iai.NewCreateGroupRequest()
	request.GroupId = common.StringPtr(groupID)
	request.GroupName = common.StringPtr(groupID)

	_, err := client.CreateGroup(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

// CreatePerson 创建人员
func CreatePerson(groupID, personID, image string) error {
	image = getImageBase64(image)

	client := getQcloudIaiClient()

	request := iai.NewCreatePersonRequest()
	request.GroupId = common.StringPtr(groupID)
	request.PersonId = common.StringPtr(personID)
	request.PersonName = common.StringPtr(personID)
	request.Image = common.StringPtr(image)

	_, err := client.CreatePerson(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

// SearchFaces 人脸搜索
func SearchFaces(groupID, image string) (string, error) {
	image = getImageBase64(image)

	client := getQcloudIaiClient()

	request := iai.NewSearchFacesRequest()
	request.GroupIds = common.StringPtrs([]string{groupID})
	request.Image = common.StringPtr(image)
	request.MaxFaceNum = common.Uint64Ptr(1)
	request.MaxPersonNum = common.Uint64Ptr(1)

	response, err := client.SearchFaces(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", err
	}
	if err != nil {
		return "", err
	}

	results := response.Response.Results
	if len(results) != 1 {
		return "", goerrors.New("no result")
	}

	candidates := results[0].Candidates
	if len(candidates) != 1 {
		return "", goerrors.New("no candidate")
	}

	return *candidates[0].PersonId, nil
}
