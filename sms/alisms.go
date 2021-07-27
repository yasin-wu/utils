package sms

import (
	"encoding/json"
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	js "github.com/bitly/go-simplejson"
)

type AliSms struct {
	Scheme          string
	RegionId        string
	AccessKeyId     string
	AccessKeySecret string
}

func New(scheme, regionId, accessKeyId, accessKeySecret string) (*AliSms, error) {
	if scheme == "" {
		scheme = "https"
	}
	if regionId == "" {
		regionId = "cn-hangzhou"
	}
	if accessKeyId == "" {
		return nil, errors.New("AccessKeyId is nil")
	}
	if accessKeySecret == "" {
		return nil, errors.New("AccessKeySecret is nil")
	}
	return &AliSms{RegionId: regionId, AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}, nil
}

func (this *AliSms) Send(signName, templateCode string, phones []string, param map[string]string) error {
	if signName == "" {
		return errors.New("SignName is nil")
	}
	if templateCode == "" {
		return errors.New("TemplateCode is nil")
	}
	if param == nil {
		return errors.New("param is nil")
	}
	phoneStr, err := this.verifyPhones(phones)
	if err != nil {
		return err
	}
	client, err := dysmsapi.NewClientWithAccessKey(this.RegionId, this.AccessKeyId, this.AccessKeySecret)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = this.Scheme
	request.PhoneNumbers = phoneStr
	request.SignName = signName
	request.TemplateCode = templateCode
	j := js.New()
	for k, v := range param {
		vr := []rune(v)
		if len(vr) > 20 {
			vr = vr[0:20]
		}
		j.Set(k, string(vr))
	}
	messageByte, err := json.Marshal(j)
	if err != nil {
		return err
	}
	request.TemplateParam = string(messageByte)
	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if response.Code != "OK" {
		return errors.New(response.Message)
	}
	return nil
}

func (this *AliSms) verifyPhones(phones []string) (string, error) {
	if phones == nil || len(phones) == 0 {
		return "", errors.New("phones is nil")
	}
	phoneStr := ""
	for _, phone := range phones {
		lenPhone := len(phone)
		if lenPhone != 11 {
			continue
		}
		phoneStr += "," + phone
	}
	if phoneStr == "" {
		return "", errors.New("invalid mobile number")
	}
	return phoneStr, nil
}
