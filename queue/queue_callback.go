package queue

import (
	"context"
	"time"

	js "github.com/bitly/go-simplejson"
	"github.com/davecgh/go-spew/spew"
	"github.com/segmentio/kafka-go"
)

func Cb(ctx context.Context, msg *kafka.Message) {
	docallback(ctx, msg, 0, nil)
}

func docallback(ctx context.Context, msg *kafka.Message, retryInterval time.Duration, nextqueue *Queue) {
	var err error
	bs := msg.Value
	j, err := js.NewJson(bs)
	if err != nil {
		return
	}
	spew.Dump(j)
}
