package consumer

import (
	"github.com/yasin-wu/utils/queue/pkg/consts"
)

type Consumer struct {
	Topics    []string      `json:"topics"`
	Name      string        `json:"name"`
	Offset    consts.Offset `json:"offset"`
	AutoAck   bool          `json:"auto_ack,default=true"`
	Exclusive bool          `json:"exclusive,default=false"`
	NoLocal   bool          `json:"no_local,default=false"`
	NoWait    bool          `json:"no_wait,default=false"`
	Reader    Reader        `json:"reader"`
}

func New(name string, topics ...string) *Consumer {
	return &Consumer{Topics: topics, Name: name, Offset: consts.OffsetNewest, Reader: &defaultReader{}}
}

func (c *Consumer) Verify() {
	if c.Offset != consts.OffsetNewest && c.Offset != consts.OffsetOldest {
		c.Offset = consts.OffsetNewest
	}
	if c.Reader == nil {
		c.Reader = &defaultReader{}
	}
}
