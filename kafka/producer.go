package kafka

import (
	"encoding/json"
	"errors"

	"github.com/Shopify/sarama"
)

type Message struct {
	Topic   string      `json:"topic"`
	Key     string      `json:"key"`
	Message interface{} `json:"message"`
}

/**
 * @author: yasinWu
 * @date: 2022/2/9 11:46
 * @params: messages []*Message
 * @return: error
 * @description: kafka producer
 */
func (k *Kafka) Send(messages []*Message) error {
	client, err := sarama.NewSyncProducer(k.brokers, k.config)
	if err != nil {
		return errors.New("new producer failed: " + err.Error())
	}
	defer client.Close()
	var msgs []*sarama.ProducerMessage
	for _, v := range messages {
		buffer, err := json.Marshal(v.Message)
		if err != nil {
			continue
		}
		msg := &sarama.ProducerMessage{
			Topic: v.Topic,
			Key:   sarama.StringEncoder(v.Key),
			Value: sarama.ByteEncoder(buffer),
		}
		msgs = append(msgs, msg)
	}
	return client.SendMessages(msgs)
}
