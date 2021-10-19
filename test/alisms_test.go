package test

import (
	"testing"

	sms2 "yasin-wu/utils/sms"
)

func TestAliSmsSend(t *testing.T) {
	sms, err := sms2.New("", "",
		"", "")
	if err != nil {
		t.Error(err)
		return
	}
	phones := []string{"18108279331"}
	param := make(map[string]string)
	param["orderno"] = "123456"
	err = sms.Send("", "SMS_185242334", phones, param)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send sms ok")
}
