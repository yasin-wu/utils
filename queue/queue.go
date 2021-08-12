package queue

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

type Callback func(ctx context.Context, msg *kafka.Message)

type Queue struct {
	Brokers         []string
	Topic           string
	Group           string
	ConsumerNum     int
	RetryNum        int
	AlreadyRetryNum int
	MinBytes        int
	MaxBytes        int
	RetryInterval   time.Duration
	Callback        Callback
	Writer          *kafka.Writer
}

func NewQueue(brokers []string, topic, group string, consumerNum int) *Queue {
	queue := Queue{}
	queue.Brokers = brokers
	queue.Topic = topic
	queue.Group = group
	queue.ConsumerNum = consumerNum
	queue.RetryNum = 1
	queue.RetryInterval = 1 * time.Minute

	queue.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.Hash{},
	})

	return &queue
}

func (this *Queue) Start() error {
	for i := 0; i < this.ConsumerNum; i++ {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  this.Brokers,
			GroupID:  this.Group,
			Topic:    this.Topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})

		go func() {
			for {
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					errMsg := fmt.Sprintf("reader.ReadMessage err:%v", err)
					fmt.Println(errMsg)
					break
				}
				this.Callback(context.Background(), &m)
			}
		}()
	}

	return nil
}

func (this *Queue) Write(key string, value []byte) error {
	err := this.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		})
	if err != nil {
		errMsg := fmt.Sprintf("Writer.WriteMessages err:%v", err)
		return errors.New(errMsg)
	}
	return nil
}

func (this *Queue) WriteBulk(key string, value []byte) error {
	if this.ConsumerNum > 1 {
		key = buildKey(key, this.ConsumerNum)
	}
	err := this.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		})
	if err != nil {
		errMsg := fmt.Sprintf("Writer.WriteMessages err:%v", err)
		return errors.New(errMsg)
	}
	return nil
}

func buildKey(Topic string, consumerNum int) string {
	randomStr := "1"
	if consumerNum > 1 {
		rand.Seed(time.Now().UTC().UnixNano())
		random := randInt(1, 1000)
		randomStr = fmt.Sprint(random)
	}
	return Topic + "?random=" + randomStr
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
