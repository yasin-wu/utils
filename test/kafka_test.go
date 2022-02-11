package test

import (
	"log"
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
	consumer()
	go producer()
	time.Sleep(30 * time.Second)
}

func producer() {
	config.Version = version
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	client := kafka.New(brokers, config)
	for i := 0; i < count; i++ {
		j := js.New()
		j.Set("num", i)
		err := client.Send(topic, key, j)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
}

func consumer() {
	config.Version = version
	client := kafka.New(brokers, config)
	client.Receive([]string{topic}, -1, nil)
}