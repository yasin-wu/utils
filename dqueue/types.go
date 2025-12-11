package dqueue

import (
	"errors"
	"time"
)

// JobAction 延迟任务执行器, 每个任务需要实现该接口, 并注册到延迟队列中, Topic唯一
type JobAction interface {
	Topic() string
	Execute(msg Message) error
}

// Message 延迟队列消息, 存放任务执行参数
type Message struct {
	ID        string         `json:"id"`         //任务ID,相同Topic下,ID唯一
	Topic     string         `json:"topic"`      //主题,JobAction的Topic()返回值,用于区分不同的任务,全局不可重复,用于注册JobAction
	Body      map[string]any `json:"body"`       //消息体,存放任务执行参数
	FaultTime int64          `json:"fault_time"` //故障重试时间,秒级,表示不重试
	Timestamp int64          `json:"timestamp"`  //创建时间戳,秒级
	ExecuteAt int64          `json:"execute_at"` //执行时间戳,秒级
}

func (msg *Message) Check() error {
	if msg.ID == "" {
		return errors.New("message id is empty")
	}
	if msg.Topic == "" {
		return errors.New("message topic is empty")
	}
	if len(msg.Body) == 0 {
		return errors.New("message body is empty")
	}
	if msg.ExecuteAt <= time.Now().Unix() {
		return errors.New("message execute_at is too old")
	}
	if msg.Timestamp <= 0 {
		msg.Timestamp = time.Now().Unix()
	}
	return nil
}
