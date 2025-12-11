package dqueue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestDQueue(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "10.34.101.49:30379",
		Password: "TopsecCdu_1130",
		DB:       15,
	})
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
			ID:    (&testDQueueJob1{}).ID(),
			Topic: "test_dqueue_topic1",
			Body: map[string]any{
				"name":       fmt.Sprintf("name1-%d", i),
				"timestamp":  time.Unix(now, 0).Format("2006-01-02 15:04:05"),
				"execute_at": time.Unix(now+int64(i+1), 0).Format("2006-01-02 15:04:05"),
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
			ID:    (&testDQueueJob2{}).ID(),
			Topic: "test_dqueue_topic2",
			Body: map[string]any{
				"name":       fmt.Sprintf("name2-%d", i),
				"timestamp":  time.Unix(now, 0).Format("2006-01-02 15:04:05"),
				"execute_at": time.Unix(now+int64(i+5), 0).Format("2006-01-02 15:04:05"),
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

type testDQueueJob1 struct{}

var _ JobAction = (*testDQueueJob1)(nil)

func (t *testDQueueJob1) ID() string {
	return "test_dqueue_job1"
}

func (t *testDQueueJob1) Execute(msg map[string]any) error {
	name, _ := msg["name"].(string)
	fmt.Printf("id: %s, name: %s, msg: %v\n", t.ID(), name, msg)
	return nil
}

type testDQueueJob2 struct{}

var _ JobAction = (*testDQueueJob2)(nil)

func (t *testDQueueJob2) ID() string {
	return "test_dqueue_job2"
}

func (t *testDQueueJob2) Execute(msg map[string]any) error {
	name, _ := msg["name"].(string)
	fmt.Printf("id: %s, name: %s, msg: %v\n", t.ID(), name, msg)
	return nil
}
