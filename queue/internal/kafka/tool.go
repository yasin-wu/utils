package kafka

import (
	"github.com/yasin-wu/utils/queue/pkg/message"

	"github.com/Shopify/sarama"
)

func (k *kafka) formatMessage(msg *sarama.ConsumerMessage) *message.ConsumerMessage {
	return &message.ConsumerMessage{
		Key:       msg.Key,
		Value:     msg.Value,
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Timestamp: msg.Timestamp,
	}
}
