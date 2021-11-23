package test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"

	queue2 "github.com/yasin-wu/utils/queue"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func TestQueue(t *testing.T) {
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConf, nil
	})
	fmt.Println("初始化Apollo配置成功")
	cache := client.GetConfigCache(apolloConf.NamespaceName)
	broker, _ := cache.Get("kafka.broker")
	brokers := strings.Split(broker.(string), ",")
	var err error
	queue := queue2.NewQueue(brokers, "yasin-test", "yasin-testGroup", 10)
	queue.Callback = queue2.Cb
	err = queue.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = queue.Write("test", []byte("test message"))
	if err != nil {
		log.Fatal(err)
	}
}
