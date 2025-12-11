package nats

import (
	"fmt"
	"strings"

	"github.com/yasin-wu/utils/queue/pkg/message"

	natsmodel "github.com/nats-io/nats.go"
)

func (n *nats) formatMessage(msg *natsmodel.Msg) *message.ConsumerMessage {
	return &message.ConsumerMessage{
		Value: msg.Data,
		Topic: msg.Subject,
	}
}

func (n *nats) toNatsMsg(messages ...*message.Message) []*natsmodel.Msg {
	var data []*natsmodel.Msg
	for _, v := range messages {
		data = append(data, &natsmodel.Msg{
			Subject: n.handleSubject(v.Topic),
			Data:    v.Message,
		})
	}
	return data
}

func (n *nats) handleSubject(topic string) string {
	subject := topic
	temp := strings.Split(topic, ".")
	if temp[0] != n.stream {
		subject = fmt.Sprintf("%s.%s", n.stream, topic)
	}
	return subject
}
