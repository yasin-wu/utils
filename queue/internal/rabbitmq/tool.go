package rabbitmq

import (
	"github.com/yasin-wu/utils/queue/pkg/message"

	"github.com/streadway/amqp"
)

func (r *rabbitMQ) formatMessage(msg amqp.Delivery) *message.ConsumerMessage {
	return &message.ConsumerMessage{
		Value:     msg.Body,
		Topic:     msg.RoutingKey,
		Timestamp: msg.Timestamp,
	}
}
