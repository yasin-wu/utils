package test

import (
	"fmt"
	"testing"
	"time"

	js "github.com/bitly/go-simplejson"

	"github.com/Shopify/sarama"
	"github.com/yasin-wu/utils/kafka"
)

var (
	brokers  = []string{"10.34.4.14:9092"}
	topic    = "test_log"
	groupId  = "test_group"
	key      = "testkey"
	strategy = "range"
	count    = 1000
	config   = kafka.NewConfig()
	version  = sarama.MaxVersion
)

func TestKafka(t *testing.T) {
	go producer()
	go consumer()
	time.Sleep(30 * time.Second)
}

func producer() {
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Version = version
	client := kafka.New(brokers, config)
	for i := 0; i < count; i++ {
		j := js.New()
		j.Set("num", i)
		pid, offset, _ := client.Send(topic, key, j)
		fmt.Printf("this is producer message, Partition:%v, Offset:%v \n", pid, offset)
		time.Sleep(time.Second)
	}
}

func consumer() {
	config.Version = version
	client := kafka.New(brokers, config)
	err := client.Receive([]string{topic}, -1)
	if err != nil {
		return
	}
}
