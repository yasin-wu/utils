package sms

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestAliSms(t *testing.T) {
	scheme := ""
	regionId := ""
	accessKeyId := ""
	accessKeySecret := ""
	phone := ""
	sms, err := New(scheme, regionId, accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatal(err)
	}
	phones := strings.Split(phone, ",")
	param := make(map[string]string)
	param["orderno"] = "123456"
	err = sms.Send("", "SMS_185242334", phones, param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("send sms ok")
}
