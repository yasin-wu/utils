package consumer

import (
	"fmt"

	"github.com/yasin-wu/utils/queue/pkg/message"
)

type Reader interface {
	Read(message ...*message.ConsumerMessage)
}

type defaultReader struct{}

var _ Reader = (*defaultReader)(nil)

func (d *defaultReader) Read(messages ...*message.ConsumerMessage) {
	for _, msg := range messages {
		fmt.Printf("this is consumer message, Topic:%s Value:%v\n",
			msg.Topic, string(msg.Value))
	}
}
