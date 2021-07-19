package test

import (
	"testing"

	js "github.com/bitly/go-simplejson"
	sms2 "github.com/yasin-wu/utils/sms"
)

func TestAliSms_Send(t *testing.T) {
	j := js.New()
	j.Set("orderno", "123456")
	sms := &sms2.AliSms{
		RegionId:        "cn-hangzhou",
		AccessKeyId:     "xxxx",
		AccessKeySecret: "xxxx",
		PhoneNumbers:    []string{"181xxxx9331"},
		SignName:        "xxxx",
		TemplateCode:    "SMS_xxxxx",
		TemplateParam:   j,
	}
	err := sms.Send()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send sms ok")
}
