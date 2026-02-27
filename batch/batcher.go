package batch

import (
	"context"
	"sync"
	"time"

	"github.com/yasin-wu/utils/util"
)

type Batcher struct {
	actuator      Actuator
	batchSize     int
	bufferSize    int
	flushInterval time.Duration
	logger        util.Logger
	dataChan      chan any
	wg            sync.WaitGroup
	mu            sync.Mutex
	batch         []any
	ctx           context.Context
	cancel        context.CancelFunc
}

type Option func(b *Batcher)

func New(actuator Actuator, opts ...Option) *Batcher {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Batcher{
		actuator:      actuator,
		batchSize:     100,
		bufferSize:    1000,
		flushInterval: 30 * time.Second,
		logger:        util.NewDefaultLogger(),
		wg:            sync.WaitGroup{},
		mu:            sync.Mutex{},
		ctx:           ctx,
		cancel:        cancel,
	}
	for _, opt := range opts {
		opt(b)
	}
	b.batch = make([]any, 0, b.batchSize)
	b.dataChan = make(chan any, b.bufferSize)
	b.wg.Add(1)
	go b.run()
	return b
}

func WithBatchSize(batchSize int) Option {
	return func(b *Batcher) {
		if batchSize > 0 {
			b.batchSize = batchSize
		}
	}
}

func WithBufferSize(bufferSize int) Option {
	return func(b *Batcher) {
		if bufferSize > 0 {
			b.bufferSize = bufferSize
		}
	}
}

func WithFlushInterval(flushInterval time.Duration) Option {
	return func(b *Batcher) {
		if flushInterval > 0 {
			b.flushInterval = flushInterval
		}
	}
}

func WithLogger(logger util.Logger) Option {
	return func(b *Batcher) {
		if logger != nil {
			b.logger = logger
		}
	}
}

func (b *Batcher) Add(data any) {
	select {
	case b.dataChan <- data:
		b.logger.Infof("add data to channel")
	case <-b.ctx.Done():
		b.logger.Errorf("batch writer has been closed")
		b.immediateWrite(data)
	}
}

func (b *Batcher) Close() error {
	b.cancel()
	b.wg.Wait()
	close(b.dataChan)
	return nil
}

func (b *Batcher) run() {
	defer b.wg.Done()
	ticker := time.NewTicker(b.flushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-b.ctx.Done():
			b.flush()
			return
		case data := <-b.dataChan:
			b.mu.Lock()
			b.batch = append(b.batch, data)
			if len(b.batch) >= b.batchSize {
				b.logger.Infof("batch flush, batch size: %d", len(b.batch))
				b.flushUnlocked()
			}
			b.mu.Unlock()
		case <-ticker.C:
			b.mu.Lock()
			if len(b.batch) > 0 {
				b.logger.Infof("writer interval flush, interval: %v", b.flushInterval)
				b.flushUnlocked()
			}
			b.mu.Unlock()
		}
	}
}

func (b *Batcher) flushUnlocked() {
	if len(b.batch) == 0 {
		return
	}
	b.logger.Infof("start flush, data count: %d", len(b.batch))
	batchToWrite := make([]any, len(b.batch))
	copy(batchToWrite, b.batch)
	b.batch = b.batch[:0]
	go func(batch []any) {
		b.bulkWrite(batch)
	}(batchToWrite)
}

func (b *Batcher) flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flushUnlocked()
}

func (b *Batcher) bulkWrite(batch []any) {
	if err := b.actuator.Bulk(batch); err != nil {
		b.logger.Errorf("bulk write failed: %v", err)
	}
}

func (b *Batcher) immediateWrite(data any) {
	if err := b.actuator.Immediate(data); err != nil {
		b.logger.Errorf("immediate write failed: %v", err)
	}
}
