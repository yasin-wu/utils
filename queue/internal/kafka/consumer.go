package kafka

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/yasin-wu/utils/queue/pkg/consumer"
)

func (k *kafka) Stop() {
	close(k.forever)
}

func (k *kafka) Topics() ([]string, error) {
	csm, err := sarama.NewConsumer(k.brokers, k.config)
	if err != nil {
		return nil, err
	}
	defer csm.Close()
	return csm.Topics()
}

func (k *kafka) Subscribe(group string, consumers ...*consumer.Consumer) {
	go func() {
		if group != "" {
			for _, csm := range consumers {
				go k.receiveGroup(group, csm)
			}
		} else {
			go k.receive(consumers...)
		}
	}()
	<-k.forever
}

func (k *kafka) receive(consumers ...*consumer.Consumer) {
	var err error
	k.ctx = context.Background()
	k.consumer, err = sarama.NewConsumer(k.brokers, k.config)
	if err != nil {
		k.logger.Errorf("new consumer failed: %v", err)
		return
	}
	defer k.consumer.Close()
	for _, csm := range consumers {
		go k.topic(csm)
	}
	select {}
}

func (k *kafka) receiveGroup(group string, consumer *consumer.Consumer) {
	keepRunning := true
	consumer.Verify()
	k.config.Consumer.Offsets.Initial = int64(consumer.Offset)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	consumerGroup, err := sarama.NewConsumerGroup(k.brokers, group, k.config)
	if err != nil {
		k.logger.Errorf("new consumer group failed: %v", err)
		return
	}
	defer consumerGroup.Close()
	csm := Consumer{
		ready:  make(chan bool),
		reader: consumer.Reader,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, consumer.Topics, &csm); err != nil {
				k.logger.Errorf("consumer failed: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			csm.ready = make(chan bool)
		}
	}()
	<-csm.ready
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	for keepRunning {
		select {
		case <-ctx.Done():
			k.logger.Errorf("terminating: context canceled")
			keepRunning = false
		case <-sigterm:
			k.logger.Errorf("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
}

func (k *kafka) topic(consumer *consumer.Consumer) {
	consumer.Verify()
	partitionList, err := k.consumer.Partitions(consumer.Topics[0])
	if err != nil {
		k.logger.Errorf("consumer partitions failed: %v", err)
		return
	}
	for partition := range partitionList {
		k.partitionConsumer, err = k.consumer.ConsumePartition(consumer.Topics[0], int32(partition), int64(consumer.Offset))
		if err != nil {
			k.logger.Errorf("start consumer for partition %d failed: %v", partition, err)
			continue
		}
		go k.message(consumer.Reader)
	}
	for {
		select {
		case <-k.ctx.Done():
			return
		default:
			continue
		}
	}
}

func (k *kafka) message(messageHandler consumer.Reader) {
	defer k.partitionConsumer.AsyncClose()
	for msg := range k.partitionConsumer.Messages() {
		messageHandler.Read(k.formatMessage(msg))
		select {
		case <-k.ctx.Done():
			return
		default:
			continue
		}
	}
}

type Consumer struct {
	ready  chan bool
	reader consumer.Reader
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		c.reader.Read((&kafka{}).formatMessage(msg))
		session.MarkMessage(msg, "")
	}
	return nil
}
