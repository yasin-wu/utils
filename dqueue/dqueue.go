package dqueue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/yasin-wu/utils/util"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type DQueue struct {
	keyPrefix  string
	batchLimit int64
	redisCli   *redis.Client
	executors  map[string]JobAction
	logger     util.Logger
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
		logger:     util.NewDefaultLogger(),
		ctx:        ctx,
		cancel:     cancel,
		wg:         sync.WaitGroup{},
		mu:         sync.RWMutex{},
	}, nil
}

func (dq *DQueue) SetLogger(logger util.Logger) {
	if logger != nil {
		dq.logger = logger
	}
}

func (dq *DQueue) Register(action JobAction) error {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	topic := action.Topic()
	if _, ok := dq.executors[topic]; ok {
		return ErrJobTopicDuplicate
	}
	dq.executors[topic] = action
	return nil
}

func (dq *DQueue) StartBackground(interval time.Duration) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
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
					go func(topic string) {
						dq.executeBatch(topic)
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
	zsetKey := dq.formatZsetKey(msg.Topic)
	hashKey := dq.formatHashKey(msg.Topic)
	member, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if len(member) == 0 {
		return errors.New("job is empty")
	}
	pipe := dq.redisCli.TxPipeline()
	oldMember, err := dq.redisCli.HGet(ctx, hashKey, msg.ID).Result()
	if err == nil && oldMember != "" {
		pipe.ZRem(ctx, zsetKey, oldMember)
	}
	z := &redis.Z{
		Score:  float64(msg.ExecuteAt),
		Member: member,
	}
	pipe.ZAdd(ctx, zsetKey, z)
	pipe.HSet(ctx, hashKey, msg.ID, string(member))
	_, err = pipe.Exec(ctx)
	return err
}

func (dq *DQueue) Remove(ctx context.Context, msg *Message) error {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	if msg.ID == "" {
		return errors.New("message id is empty")
	}
	if msg.Topic == "" {
		return errors.New("message topic is empty")
	}
	member, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if len(member) == 0 {
		return errors.New("job is empty")
	}
	zsetKey := dq.formatZsetKey(msg.Topic)
	hashKey := dq.formatHashKey(msg.Topic)
	pipe := dq.redisCli.TxPipeline()
	oldMember, err := dq.redisCli.HGet(ctx, hashKey, msg.ID).Result()
	if err == nil && oldMember != "" {
		pipe.ZRem(ctx, zsetKey, oldMember)
	}
	pipe.HDel(ctx, hashKey, msg.ID)
	pipe.ZRem(ctx, zsetKey, member)
	_, err = pipe.Exec(ctx)
	return err
}

func (dq *DQueue) executeBatch(topic string) {
	messages, err := dq.getReadyMessages(dq.ctx, topic)
	if err != nil {
		dq.logger.Errorf("get ready messages error: %v", err)
		return
	}
	for _, msg := range messages {
		executor, ok := dq.executors[msg.Topic]
		if !ok {
			continue
		}
		if err := executor.Execute(msg); err != nil {
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

func (dq *DQueue) getReadyMessages(ctx context.Context, topic string) ([]Message, error) {
	now := time.Now().Unix()
	opt := &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%d", now),
		Offset: 0,
		Count:  dq.batchLimit,
	}
	key := dq.formatZsetKey(topic)
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

func (dq *DQueue) formatZsetKey(topic string) string {
	return fmt.Sprintf("%s:%s", dq.keyPrefix, topic)
}

func (dq *DQueue) formatHashKey(topic string) string {
	return fmt.Sprintf("%s:%s:id_index", dq.keyPrefix, topic)
}
