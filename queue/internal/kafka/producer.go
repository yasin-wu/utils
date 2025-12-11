package kafka

import (
	"errors"

	"github.com/yasin-wu/utils/queue/pkg/message"

	"github.com/Shopify/sarama"
)

func (k *kafka) Publish(messages ...*message.Message) error {
	client, err := sarama.NewSyncProducer(k.brokers, k.config)
	if err != nil {
		return errors.New("new producer failed: " + err.Error())
	}
	defer client.Close()
	var msgs []*sarama.ProducerMessage
	for _, v := range messages {
		msg := &sarama.ProducerMessage{
			Topic: v.Topic,
			Key:   sarama.StringEncoder(v.Key),
			Value: sarama.ByteEncoder(v.Message),
		}
		msgs = append(msgs, msg)
	}
	return client.SendMessages(msgs)
}
