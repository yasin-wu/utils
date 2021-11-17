package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"

	email2 "github.com/yasin-wu/utils/email"
)

func TestEmailSend(t *testing.T) {
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConf, nil
	})
	fmt.Println("初始化Apollo配置成功")
	cache := client.GetConfigCache(apolloConf.NamespaceName)
	host, _ := cache.Get("email.host")
	port, _ := cache.Get("email.port")
	user, _ := cache.Get("email.user")
	password, _ := cache.Get("email.password")
	from, _ := cache.Get("email.from")
	to, _ := cache.Get("email.to")
	email, err := email2.New(host.(string), port.(string), user.(string), password.(string), from.(string))
	if err != nil {
		t.Error(err)
		return
	}
	tos := strings.Split(to.(string), ",")
	err = email.SendTLS(tos, "test", "test")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("send email ok")
}
