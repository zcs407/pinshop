package controllers

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func MsgSend(phoneNmber, vscode string) error {
	fmt.Println(phoneNmber, vscode)
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "LTAI49yQmf3Tbhdi", "dDNrUp9tKQK4kOORDXMNIkWV23dl4R")
	if err != nil {
		fmt.Println("NewClientWithAccessKey err:", err)
		return err
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phoneNmber
	request.QueryParams["SignName"] = "品优购"
	request.QueryParams["TemplateCode"] = "SMS_164275022"
	request.QueryParams["TemplateParam"] = `{"code":` + vscode + `}`
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Println("resonse err", err)
		return err
	}
	fmt.Print(response.GetHttpContentString())
	return nil
}
