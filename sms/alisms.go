package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	js "github.com/bitly/go-simplejson"
	"github.com/yasin-wu/utils/common"
)

/**
 * @author: yasin
 * @date: 2020/10/23 09:11
 * @description：TemplateParam's key is template set
 */
type AliSms struct {
	RegionId        string   `json:"region_id"`
	AccessKeyId     string   `json:"access_key_id"`
	AccessKeySecret string   `json:"access_key_secret"`
	PhoneNumbers    []string `json:"phone_numbers"`
	SignName        string   `json:"sign_name"`
	TemplateCode    string   `json:"template_code"`
	TemplateParam   *js.Json `json:"template_param"`
}

func (this *AliSms) Send() error {
	err := this.verifyRequired()
	if err != nil {
		return err
	}
	phones, err := this.verifyPhoneNumbers()
	if err != nil {
		return err
	}

	client, err := dysmsapi.NewClientWithAccessKey(
		this.RegionId,
		this.AccessKeyId,
		this.AccessKeySecret)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = common.Ali_SMS_SCHEME
	request.PhoneNumbers = phones
	request.SignName = this.SignName
	request.TemplateCode = this.TemplateCode
	if this.TemplateParam != nil {
		//参数不能超过20字符
		paramMap, _ := this.TemplateParam.Map()
		paramJson := js.New()
		for k, v := range paramMap {
			sv := fmt.Sprintf("%v", v)
			if len(sv) > 20 {
				sv = sv[0:18]
			}
			paramJson.Set(k, sv)
		}
		messageByte, err := json.Marshal(paramJson)
		if err != nil {
			return err
		}
		messageStr := string(messageByte)
		request.TemplateParam = messageStr
	}
	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if response.Code != common.Ali_SMS_SUCCESS {
		return errors.New(response.Message)
	}
	return nil
}

func (this *AliSms) verifyPhoneNumbers() (string, error) {
	if this.PhoneNumbers == nil || len(this.PhoneNumbers) == 0 {
		return "", errors.New("PhoneNumbers is nil")
	}
	phones := ""
	for _, phone := range this.PhoneNumbers {
		lenPhone := len(phone)
		if lenPhone != 11 {
			fmt.Println(phone + "invalid mobile number")
			continue
		}
		phones += "," + phone
	}
	if phones == "" {
		return "", errors.New("invalid mobile number")
	}
	return phones, nil
}

func (this *AliSms) verifyRequired() error {
	if this.RegionId == "" {
		return errors.New("RegionId is nil")
	}
	if this.AccessKeyId == "" {
		return errors.New("AccessKeyId is nil")
	}
	if this.AccessKeySecret == "" {
		return errors.New("AccessKeySecret is nil")
	}
	if this.SignName == "" {
		return errors.New("SignName is nil")
	}
	if this.TemplateCode == "" {
		return errors.New("TemplateCode is nil")
	}
	return nil
}
