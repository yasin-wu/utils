package test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"

	sms2 "github.com/yasin-wu/utils/sms"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func TestAliSmsSend(t *testing.T) {
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConf, nil
	})
	fmt.Println("初始化Apollo配置成功")
	cache := client.GetConfigCache(apolloConf.NamespaceName)
	scheme, _ := cache.Get("sms.scheme")
	regionId, _ := cache.Get("sms.region_id")
	accessKeyId, _ := cache.Get("sms.access_key_id")
	accessKeySecret, _ := cache.Get("sms.access_key_secret")
	phone, _ := cache.Get("sms.to")
	sms, err := sms2.New(scheme.(string), regionId.(string), accessKeyId.(string), accessKeySecret.(string))
	if err != nil {
		log.Fatal(err)
	}
	phones := strings.Split(phone.(string), ",")
	param := make(map[string]string)
	param["orderno"] = "123456"
	err = sms.Send("", "SMS_185242334", phones, param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("send sms ok")
}
