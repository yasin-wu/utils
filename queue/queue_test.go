package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	var brokers = []string{"127.0.0.1:9092"}
	var err error
	queue := NewQueue(brokers, "yasin-test", "yasin-testGroup", 10)
	queue.Callback = cb
	err = queue.Start()
	if err != nil {
		t.Log(err)
		panic(err)
	}
	err = queue.Write("test", []byte("test message"))
	if err != nil {
		t.Log(err)
	}
}
