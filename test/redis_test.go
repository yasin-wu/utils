package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"

	"github.com/yasin-wu/utils/redis"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

var key = "test-redis"

func TestRedis(t *testing.T) {
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConf, nil
	})
	fmt.Println("初始化Apollo配置成功")
	cache := client.GetConfigCache(apolloConf.NamespaceName)
	host, _ := cache.Get("redis.host")
	password, _ := cache.Get("redis.password")
	cli, err := redis.New(host.(string), redis.WithPassWord(password.(string)))
	if err != nil {
		log.Fatal(err)
	}
	cli.Set(key, "", time.Minute)
	fmt.Println(cli.TTL(key))
}
