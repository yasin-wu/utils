package kafka

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

func (k *Kafka) Receive(topics []string, offset int64) error {
	var err error
	if k.groupId != "" && k.strategy != "" {
		err = k.receiveGroup(topics, offset)
	} else {
		err = k.receive(topics, offset)
	}
	return err
}

func (k *Kafka) receive(topics []string, offset int64) error {
	var err error
	k.config.Consumer.Return.Errors = true
	k.consumer, err = sarama.NewConsumer(k.brokers, k.config)
	if err != nil {
		return errors.New("new consumer failed: " + err.Error())
	}
	defer k.consumer.Close()
	for _, v := range topics {
		k.ctx = context.Background()
		go k.topic(v, offset)
	}
	select {}
}

func (k *Kafka) receiveGroup(topics []string, offset int64) error {
	keepRunning := true
	k.config.Consumer.Return.Errors = true
	k.config.Consumer.Offsets.Initial = offset
	switch k.strategy {
	case "sticky":
		k.config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		k.config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		k.config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		return errors.New("strategy is not supported")
	}
	ctx, cancel := context.WithCancel(context.Background())
	consumerGroup, err := sarama.NewConsumerGroup(k.brokers, k.groupId, k.config)
	if err != nil {
		log.Fatalf("new consumer group failed: %v", err)
		return err
	}
	defer consumerGroup.Close()
	consumer := Consumer{
		ready:          make(chan bool),
		messageHandler: k.messageHandler,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, topics, &consumer); err != nil {
				log.Fatalf("error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	return nil
}

func (k *Kafka) topic(topic string, offset int64) {
	partitionList, err := k.consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("fail to start consumer partition,err:%v\n", err)
		return
	}
	for partition := range partitionList {
		k.partitionConsumer, err = k.consumer.ConsumePartition(topic, int32(partition), offset)
		if err != nil {
			log.Printf("fail to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		go k.message()
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

func (k *Kafka) message() {
	defer k.partitionConsumer.AsyncClose()
	for message := range k.partitionConsumer.Messages() {
		k.messageHandler(message)
		select {
		case <-k.ctx.Done():
			return
		default:
			continue
		}
	}
}

type Consumer struct {
	ready          chan bool
	messageHandler MessageHandler
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.messageHandler(message)
		session.MarkMessage(message, "")
	}
	return nil
}
