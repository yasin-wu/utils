package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/yasin-wu/utils/queue/internal/kafka"
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/queue/pkg/consts"
	"github.com/yasin-wu/utils/queue/pkg/consumer"
	"github.com/yasin-wu/utils/queue/pkg/message"
)

var queueConfig = &Config{
	Brokers:     []string{"10.34.101.32:30992"},
	Driver:      "kafka",
	Username:    "guest",
	Password:    "guest",
	KafkaConfig: kafka.NewConfig(),
	NatsConfig: &config.NatsConfig{
		Stream:        "yasin",
		StreamEnabled: true,
		NCOpt:         nil,
		JSOpt:         nil,
	},
}

var (
	topics = []string{"test.messages"}
	count  = 10
)

func TestConsumer(t *testing.T) {
	//queueConfig.KafkaConfig.Version = sarama.V2_1_1_0
	cli, err := New(queueConfig)
	if err != nil {
		log.Fatal(err)
	}
	var csms []*consumer.Consumer
	for _, v := range topics {
		csm := consumer.New(v, v)
		csm.Offset = consts.OffsetOldest
		csm.Reader = &consumerReader{}
		csms = append(csms, csm)
	}
	cli.Subscribe("", csms...)
	runtime.Goexit()
}

func TestProducer(t *testing.T) {
	queueClient, err := New(queueConfig)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range topics {
		var msgs []*message.Message
		for i := 0; i < count; i++ {
			j := make(map[string]interface{})
			j["id"] = i
			j["timestamp"] = time.Now().Unix()
			buffer, _ := json.Marshal(j)
			msg := &message.Message{
				Topic:   v,
				Message: buffer,
			}
			msgs = append(msgs, msg)
		}
		err = queueClient.Publish(msgs...)
		if err != nil {
			log.Println(err)
		}
	}
}

type consumerReader struct{}

func (c *consumerReader) Read(messages ...*message.ConsumerMessage) {
	for _, msg := range messages {
		fmt.Printf("Topic:%s Value:%v\n", msg.Topic, string(msg.Value))
	}
}
