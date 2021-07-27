package test

import (
	"testing"

	sms2 "github.com/yasin-wu/utils/sms"
)

func TestAliSmsSend(t *testing.T) {
	sms, err := sms2.New("", "",
		"xxxx", "xxxx")
	if err != nil {
		t.Error(err)
		return
	}
	phones := []string{"xxxx"}
	param := make(map[string]string)
	param["xxxx"] = "123456"
	err = sms.Send("xxxx", "SMS_xxxx", phones, param)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send sms ok")
}
