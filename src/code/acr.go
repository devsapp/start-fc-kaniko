package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/google/uuid"
)

type AcrGetAuthTokenData struct {
	Data struct {
		AuthorizationToken string `json:"authorizationToken"`
		TempUserName       string `json:"tempUserName"`
		ExpireDate         int64  `json:"expireDate"`
	} `json:"data"`
}

type ServerErrorData struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func CreateAcrClient(accessKeyId *string, accessKeySecret *string, securityToken *string, region string) (_result *sdk.Client, _err error) {
	if securityToken == nil {
		_result, _err = sdk.NewClientWithAccessKey(region, *accessKeyId, *accessKeySecret)
	} else {
		_result, _err = sdk.NewClientWithStsToken(region, *accessKeyId, *accessKeySecret, *securityToken)
	}
	return
}

func GetAcrTmpUser(creds AliyunCredentials, region string) (*AcrGetAuthTokenData, error) {
	client, err := CreateAcrClient(creds.AccessKeyId, creds.AccessKeySecret, creds.SecurityToken, region)
	if err != nil {
		return nil, err
	}
	resp, err := getAuthTokenData(client, region)
	if err != nil {
		if se, ok := err.(*errors.ServerError); ok {
			fmt.Println(se.ErrorCode(), se.Message())
			var serverErrorData ServerErrorData
			err = json.Unmarshal([]byte(se.Message()), &serverErrorData)
			if err == nil && serverErrorData.Code == "USER_NOT_EXIST" {
				_err := createUserInfo(client, region)
				if _err == nil { // retry once
					return getAuthTokenData(client, region)
				}
			}
		}
	}
	return resp, err
}

func getAuthTokenData(client *sdk.Client, region string) (*AcrGetAuthTokenData, error) {
	request := requests.NewCommonRequest()
	request.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", region)
	request.Version = "2016-06-07"
	request.PathPattern = "/tokens"
	request.Method = "GET"
	request.SetContentType("application/json")
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	// fmt.Print(response.GetHttpContentString())
	var authTokenData AcrGetAuthTokenData
	err = json.Unmarshal(response.GetHttpContentBytes(), &authTokenData)
	if err != nil {
		return nil, err
	}
	fmt.Println(authTokenData.Data.TempUserName, authTokenData.Data.AuthorizationToken)
	return &authTokenData, nil
}

func createUserInfo(client *sdk.Client, region string) error {
	request := requests.NewCommonRequest()
	request.Domain = fmt.Sprintf("cr.%s.aliyuncs.com", region)
	request.Version = "2016-06-07"
	request.PathPattern = "/users"
	request.Method = "PUT"
	request.SetContentType("application/json")
	pwd := "Pwd_" + uuid.New().String()[:6]
	request.Content = []byte(fmt.Sprintf(`{"User":{"Password":"%s"}}`, pwd))
	_, err := client.ProcessCommonRequest(request)
	return err
}
