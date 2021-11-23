package queue

import (
	"context"
	"encoding/json"
	"fmt"

	js "github.com/bitly/go-simplejson"
	"github.com/segmentio/kafka-go"
)

func Cb(ctx context.Context, msg *kafka.Message) {
	doCallback(ctx, msg)
}

func doCallback(ctx context.Context, msg *kafka.Message) {
	var err error
	bs := msg.Value
	j, err := js.NewJson(bs)
	if err != nil {
		return
	}
	jstr, _ := json.MarshalIndent(j, "", "\t")
	fmt.Println(string(jstr))
}
