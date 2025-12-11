package dqueue

import (
	"errors"
	"time"
)

// JobAction 延迟任务执行器, 每个任务需要实现该接口, 并注册到延迟队列中, ID唯一
type JobAction interface {
	ID() string
	Execute(msg map[string]any) error
}

// Message 延迟队列消息, 存放任务执行参数
type Message struct {
	ID        string         `json:"id"`         //任务ID, JobAction的ID
	Topic     string         `json:"topic"`      //主题
	Body      map[string]any `json:"body"`       //消息体, 存放任务执行参数
	Timestamp int64          `json:"timestamp"`  //创建时间戳,秒级
	ExecuteAt int64          `json:"execute_at"` //执行时间戳,秒级
}

func (msg *Message) Check() error {
	if msg.ID == "" {
		return errors.New("message id is empty")
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
