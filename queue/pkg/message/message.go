package message

import (
	"time"
)

type Message struct {
	Topic   string `json:"topic"`
	Key     string `json:"key"`
	Message []byte `json:"message"`
}

type ConsumerMessage struct {
	Key, Value []byte
	Topic      string
	Partition  int32
	Offset     int64
	Timestamp  time.Time
}
