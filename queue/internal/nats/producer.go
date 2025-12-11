package nats

import (
	"errors"
	"strings"

	"github.com/yasin-wu/utils/queue/pkg/message"
)

func (n *nats) Publish(messages ...*message.Message) error {
	msg := n.toNatsMsg(messages...)
	if len(msg) == 0 {
		return errors.New("message is nil")
	}
	var errMsg []string
	for _, v := range msg {
		var err error
		switch {
		case !n.streamEnabled:
			err = n.conn.PublishMsg(v)
		case n.streamEnabled && n.async:
			_, err = n.jetStream.PublishMsgAsync(v)
		case n.streamEnabled && !n.async:
			_, err = n.jetStream.PublishMsg(v)
		}
		if err != nil {
			errMsg = append(errMsg, err.Error())
		}
	}
	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}
	return nil
}
