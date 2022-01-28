package kafka

import (
	"encoding/json"
	"errors"

	"github.com/Shopify/sarama"
)

func (k *Kafka) Send(topic, key string, value interface{}) (int32, int64, error) {
	k.config.Producer.Return.Successes = true
	client, err := sarama.NewSyncProducer(k.brokers, k.config)
	if err != nil {
		return 0, 0, errors.New("new producer failed: " + err.Error())
	}
	defer client.Close()
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Key = sarama.StringEncoder(key)
	buffer, err := json.Marshal(value)
	if err != nil {
		return 0, 0, errors.New("json marshal error: " + err.Error())
	}
	msg.Value = sarama.ByteEncoder(buffer)
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		return 0, 0, errors.New("send message failed: " + err.Error())
	}
	return pid, offset, nil
}
