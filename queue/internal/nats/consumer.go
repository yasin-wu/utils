package nats

import (
	natsmodel "github.com/nats-io/nats.go"

	"github.com/yasin-wu/utils/queue/pkg/consts"
	"github.com/yasin-wu/utils/queue/pkg/consumer"
)

func (n *nats) Stop() {
	close(n.forever)
}

func (n *nats) Subscribe(group string, consumers ...*consumer.Consumer) {
	go func() {
		for _, v := range consumers {
			go n.subscribe(group, v)
		}
	}()
	<-n.forever
}

func (n *nats) subscribe(group string, consumer *consumer.Consumer) {
	consumer.Verify()
	csm := &Consumer{
		reader: consumer.Reader,
	}
	subject := n.handleSubject(consumer.Topics[0])
	if n.streamEnabled {
		n.streamSubscribe(subject, csm, consumer)
	} else {
		var err error
		if group != "" {
			_, err = n.conn.QueueSubscribe(subject, group, csm.natsHandler)
		} else {
			_, err = n.conn.Subscribe(subject, csm.natsHandler)
		}
		if err != nil {
			n.logger.Errorf("subscribe failed : %v, subject : %s", err, subject)
		}
	}
}

func (n *nats) streamSubscribe(subject string, csm *Consumer, consumer *consumer.Consumer) {
	var subOpt []natsmodel.SubOpt
	if consumer.Offset == consts.OffsetNewest {
		subOpt = append(subOpt, natsmodel.DeliverNew())
	} else if consumer.Offset == consts.OffsetOldest {
		subOpt = append(subOpt, natsmodel.DeliverAll())
	}
	if consumer.Name != "" {
		if err := n.jetStream.DeleteConsumer(n.stream, consumer.Name); err == nil {
			subOpt = append(subOpt, natsmodel.Durable(consumer.Name))
		}
	}
	_, err := n.jetStream.Subscribe(subject, csm.natsHandler, subOpt...)
	if err != nil {
		n.logger.Errorf("stream subscribe failed : %v, subject : %s", err, subject)
	}
}

type Consumer struct {
	reader consumer.Reader
}

func (c *Consumer) natsHandler(msg *natsmodel.Msg) {
	c.reader.Read((&nats{}).formatMessage(msg))
}
