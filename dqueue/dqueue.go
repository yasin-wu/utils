package dqueue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type DQueue struct {
	keyPrefix  string
	batchLimit int64
	redisCli   *redis.Client
	executors  map[string]JobAction
	logger     Logger
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
}

func New(keyPrefix string, batchLimit int64, redisCli *redis.Client) (*DQueue, error) {
	err := redisCli.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &DQueue{
		keyPrefix:  keyPrefix,
		batchLimit: batchLimit,
		redisCli:   redisCli,
		executors:  make(map[string]JobAction),
		logger:     DefaultLogger,
		ctx:        ctx,
		cancel:     cancel,
		wg:         sync.WaitGroup{},
		mu:         sync.RWMutex{},
	}, nil
}

func (dq *DQueue) SetLogger(logger Logger) {
	if logger != nil {
		dq.logger = logger
	}
}

func (dq *DQueue) Register(action JobAction) error {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	id := action.ID()
	if _, ok := dq.executors[id]; ok {
		return ErrJobIDDuplicate
	}
	dq.executors[id] = action
	return nil
}

func (dq *DQueue) StartBackground(interval time.Duration) {
	dq.mu.Lock()
	dq.mu.Unlock()
	dq.wg.Add(1)
	go func() {
		defer dq.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-dq.ctx.Done():
				dq.logger.Infof("dqueue goroutine shutting down")
				return
			case <-ticker.C:
				for k := range dq.executors {
					go func(id string) {
						dq.executeBatch(id)
					}(k)
				}
			}
		}
	}()
}

func (dq *DQueue) Stop() {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	dq.cancel()
	dq.wg.Wait()
	dq.logger.Infof("dqueue stopped")
}

func (dq *DQueue) Add(ctx context.Context, msg *Message) error {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	if err := msg.Check(); err != nil {
		return err
	}
	key := dq.formatKey(msg.ID)
	member, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if len(member) == 0 {
		return errors.New("job is empty")
	}
	var z redis.Z
	z.Member = member
	z.Score = float64(msg.ExecuteAt)
	return dq.redisCli.ZAdd(ctx, key, &z).Err()
}

func (dq *DQueue) Remove(ctx context.Context, msg *Message) error {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	if msg.ID == "" {
		return errors.New("message id is empty")
	}
	key := dq.formatKey(msg.ID)
	member, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if len(member) == 0 {
		return errors.New("job is empty")
	}
	return dq.redisCli.ZRem(ctx, key, member).Err()
}

func (dq *DQueue) executeBatch(id string) {
	messages, err := dq.getReadyMessages(dq.ctx, id)
	if err != nil {
		dq.logger.Errorf("get ready messages error: %v", err)
		return
	}
	for _, msg := range messages {
		executor, ok := dq.executors[msg.ID]
		if !ok {
			continue
		}
		if err := executor.Execute(msg.Body); err != nil {
			dq.logger.Errorf("job action execute failed, error:%v", err)
			if msg.FaultTime <= 0 {
				continue
			}
			msg.ExecuteAt += msg.FaultTime
			if err := dq.Add(dq.ctx, &msg); err != nil {
				dq.logger.Errorf("add message error: %v", err)
			}
			continue
		}
		if err := dq.Remove(dq.ctx, &msg); err != nil {
			dq.logger.Errorf("remove message error: %v", err)
			continue
		}
	}
}

func (dq *DQueue) getReadyMessages(ctx context.Context, id string) ([]Message, error) {
	now := time.Now().Unix()
	opt := &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%d", now),
		Offset: 0,
		Count:  dq.batchLimit,
	}
	key := dq.formatKey(id)
	members, err := dq.redisCli.ZRangeByScore(ctx, key, opt).Result()
	if err != nil {
		return nil, fmt.Errorf("get ready messages error: %v", err)
	}
	var messages []Message
	for _, member := range members {
		var msg Message
		if err := json.Unmarshal([]byte(member), &msg); err != nil {
			dq.logger.Errorf("unmarshal member error: %v", err)
			continue
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (dq *DQueue) formatKey(name string) string {
	return fmt.Sprintf("%s:%s", dq.keyPrefix, name)
}
