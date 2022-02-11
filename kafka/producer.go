package kafka

import (
	"encoding/json"
	"errors"

	"github.com/Shopify/sarama"
)

/**
 * @author: yasinWu
 * @date: 2022/2/9 11:46
 * @params: topic, key string, value interface{}
 * @return: error
 * @description: kafka producer
 */
func (k *Kafka) Send(topic, key string, value interface{}) error {
	buffer, err := json.Marshal(value)
	if err != nil {
		return errors.New("json marshal error: " + err.Error())
	}
	k.config.Producer.Return.Successes = true
	client, err := sarama.NewSyncProducer(k.brokers, k.config)
	if err != nil {
		return errors.New("new producer failed: " + err.Error())
	}
	defer client.Close()
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Key = sarama.StringEncoder(key)
	msg.Value = sarama.ByteEncoder(buffer)
	_, _, err = client.SendMessage(msg)
	return err
}
