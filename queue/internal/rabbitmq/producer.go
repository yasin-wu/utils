package rabbitmq

import (
	"errors"
	"strings"

	"github.com/yasin-wu/utils/queue/pkg/message"

	"github.com/streadway/amqp"
)

func (r *rabbitMQ) Publish(messages ...*message.Message) error {
	var errMsg []string
	for _, m := range messages {
		err := r.channel.Publish(
			r.exchange,
			m.Topic,
			r.mandatory,
			r.immediate,
			amqp.Publishing{
				ContentType: r.contentType,
				Body:        m.Message,
			},
		)
		if err != nil {
			errMsg = append(errMsg, err.Error())
			continue
		}
	}
	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ";"))
	}
	return nil
}
