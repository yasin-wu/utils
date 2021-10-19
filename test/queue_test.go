package test

import (
	"testing"

	queue2 "github.com/yasin-wu/utils/queue"
)

func TestQueue(t *testing.T) {
	var brokers = []string{"127.0.0.1:9092"}
	var err error
	queue := queue2.NewQueue(brokers, "yasin-test", "yasin-testGroup", 10)
	queue.Callback = queue2.Cb
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
