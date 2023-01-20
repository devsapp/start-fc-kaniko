package main

import (
	"fmt"
	cr20181201 "github.com/alibabacloud-go/cr-20181201/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func CreateAcrEEClient(accessKeyId *string, accessKeySecret *string, securityToken *string, region string) (_result *cr20181201.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		SecurityToken:   securityToken,
	}
	config.Endpoint = tea.String(fmt.Sprintf("cr.%s.aliyuncs.com", region))
	_result = &cr20181201.Client{}
	_result, _err = cr20181201.NewClient(config)
	return _result, _err
}

func GetAcrEETmpUser(creds AliyunCredentials, region string, instanceId string) (*cr20181201.GetAuthorizationTokenResponseBody, error) {
	client, _err := CreateAcrEEClient(creds.AccessKeyId, creds.AccessKeySecret, creds.SecurityToken, region)
	if _err != nil {
		return nil, _err
	}
	getAuthorizationTokenRequest := &cr20181201.GetAuthorizationTokenRequest{InstanceId: tea.String(instanceId)}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.GetAuthorizationTokenWithOptions(getAuthorizationTokenRequest, runtime)
	if _err != nil {
		return nil, _err
	}
	return resp.Body, nil
}
