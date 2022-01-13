package test

import (
	"log"
	"strings"
	"testing"

	queue2 "github.com/yasin-wu/utils/queue"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func TestQueue(t *testing.T) {
	broker := "127.0.0.1:9092"
	brokers := strings.Split(broker, ",")
	var err error
	queue := queue2.NewQueue(brokers, "yasin-test", "yasin-testGroup", 10)
	queue.Callback = queue2.Cb
	err = queue.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = queue.Write("test", []byte("test message"))
	if err != nil {
		log.Fatal(err)
	}
}
