package sms

import (
	"encoding/json"
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	js "github.com/bitly/go-simplejson"
)

type AliSms struct {
	scheme          string
	regionID        string
	accessKeyID     string
	accessKeySecret string
}

func New(scheme, regionID, accessKeyID, accessKeySecret string) (*AliSms, error) {
	if scheme == "" {
		scheme = "https"
	}
	if regionID == "" {
		regionID = "cn-hangzhou"
	}
	if accessKeyID == "" {
		return nil, errors.New("AccessKeyId is nil")
	}
	if accessKeySecret == "" {
		return nil, errors.New("AccessKeySecret is nil")
	}
	return &AliSms{scheme: scheme, regionID: regionID, accessKeyID: accessKeyID, accessKeySecret: accessKeySecret}, nil
}

func (a *AliSms) Send(signName, templateCode string, phones []string, param map[string]string) error {
	if signName == "" {
		return errors.New("SignName is nil")
	}
	if templateCode == "" {
		return errors.New("TemplateCode is nil")
	}
	if param == nil {
		return errors.New("param is nil")
	}
	phoneStr, err := a.verifyPhones(phones)
	if err != nil {
		return err
	}
	client, err := dysmsapi.NewClientWithAccessKey(a.regionID, a.accessKeyID, a.accessKeySecret)
	if err != nil {
		return err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = a.scheme
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

func (a *AliSms) verifyPhones(phones []string) (string, error) {
	if len(phones) == 0 {
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
