package rabbitmq

import (
	"github.com/yasin-wu/utils/queue/pkg/consumer"

	"github.com/streadway/amqp"
)

func (r *rabbitMQ) Stop() {
	r.channel.Close()
	close(r.forever)
}

func (r *rabbitMQ) Topics() ([]string, error) {
	return nil, nil
}

func (r *rabbitMQ) Subscribe(_ string, consumers ...*consumer.Consumer) {
	for _, csm := range consumers {
		csm.Verify()
		if err := r.exchangeDeclare(); err != nil {
			r.logger.Errorf("exchange declare failed: %s", err)
			return
		}
		queue, err := r.queueDeclare(csm.Topics[0])
		if err != nil {
			r.logger.Errorf("queue declare failed: %s", err)
			return
		}
		if err = r.queueBind(queue.Name, csm.Topics[0]); err != nil {
			r.logger.Errorf("queue bind failed: %s", err)
			return
		}
		msg, err := r.consume(csm)
		if err != nil {
			r.logger.Errorf("listener failed: %v", err)
			return
		}
		go func(consumer *consumer.Consumer) {
			for d := range msg {
				consumer.Reader.Read(r.formatMessage(d))
			}
		}(csm)
	}
	<-r.forever
}

func (r *rabbitMQ) exchangeDeclare() error {
	return r.channel.ExchangeDeclare(
		r.exchange,
		r.direct,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *rabbitMQ) queueDeclare(topic string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *rabbitMQ) queueBind(name, topic string) error {
	return r.channel.QueueBind(
		name,
		topic,
		r.exchange,
		false,
		nil,
	)
}

func (r *rabbitMQ) consume(consumer *consumer.Consumer) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		consumer.Topics[0],
		consumer.Name,
		consumer.AutoAck,
		consumer.Exclusive,
		consumer.NoLocal,
		consumer.NoWait,
		nil,
	)
}
