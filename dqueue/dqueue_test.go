package dqueue

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisCli = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       15,
})

func TestDQueue(t *testing.T) {
	dq, err := New("test_dqueue", 10000, redisCli)
	if err != nil {
		t.Fatal(err)
	}
	err = dq.Register(&testDQueueJob1{})
	err = dq.Register(&testDQueueJob2{})
	if err != nil {
		t.Fatal(err)
	}
	dq.StartBackground(time.Second)
	for i := 0; i < 10; i++ {
		now := time.Now().Unix()
		err = dq.Add(context.Background(), &Message{
			ID:    strconv.Itoa(i),
			Topic: (&testDQueueJob1{}).Topic(),
			Body: map[string]any{
				"name": fmt.Sprintf("name1-%d", i),
			},
			Timestamp: now,
			ExecuteAt: now + int64(i+1),
		})
		if err != nil {
			t.Error()
			continue
		}
	}
	for i := 0; i < 10; i++ {
		now := time.Now().Unix()
		err = dq.Add(context.Background(), &Message{
			ID:    fmt.Sprintf("%d", i),
			Topic: (&testDQueueJob2{}).Topic(),
			Body: map[string]any{
				"name": fmt.Sprintf("name2-%d", i),
			},
			Timestamp: now,
			ExecuteAt: now + int64(i+5),
		})
		if err != nil {
			t.Error()
			continue
		}
	}
	time.Sleep(1 * time.Minute)
}

func TestAdd(t *testing.T) {
	dq, err := New("test_dqueue", 10000, redisCli)
	if err != nil {
		t.Fatal(err)
	}
	err = dq.Add(context.Background(), &Message{
		ID:    "1",
		Topic: "test_add",
		Body: map[string]any{
			"name": fmt.Sprintf("name1-%d", 1),
		},
		Timestamp: time.Now().Unix(),
		ExecuteAt: time.Now().Unix() + 10,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemove(t *testing.T) {
	dq, err := New("test_dqueue", 10000, redisCli)
	if err != nil {
		t.Fatal(err)
	}
	err = dq.Remove(context.Background(), &Message{
		ID:    "1",
		Topic: "test_add",
	})
	if err != nil {
		t.Fatal(err)
	}
}

type testDQueueJob1 struct{}

var _ JobAction = (*testDQueueJob1)(nil)

func (t *testDQueueJob1) Topic() string {
	return "test_dqueue_job1"
}

func (t *testDQueueJob1) Execute(msg Message) error {
	fmt.Printf("topic: %s, id: %s, timestamp: %s, execute_at: %s, msg: %v\n",
		msg.Topic, msg.ID, time.Unix(msg.Timestamp, 0), time.Unix(msg.ExecuteAt, 0), msg)
	return nil
}

type testDQueueJob2 struct{}

var _ JobAction = (*testDQueueJob2)(nil)

func (t *testDQueueJob2) Topic() string {
	return "test_dqueue_job2"
}

func (t *testDQueueJob2) Execute(msg Message) error {
	fmt.Printf("topic: %s, id: %s, timestamp: %s, execute_at: %s, msg: %v\n",
		msg.Topic, msg.ID, time.Unix(msg.Timestamp, 0), time.Unix(msg.ExecuteAt, 0), msg)
	return nil
}
