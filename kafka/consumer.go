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

type MessageHandler func(message *sarama.ConsumerMessage)

/**
 * @author: yasinWu
 * @date: 2022/2/9 11:45
 * @params: topics []string, offset int64
 * @return: error
 * @description: kafka consumer
 */
func (k *Kafka) Receive(topics []string, offset int64, messageHandler MessageHandler) {
	if messageHandler == nil {
		messageHandler = printMsg
	}
	go func() {
		var err error
		if k.groupId != "" && k.strategy != "" {
			err = k.receiveGroup(topics, offset, messageHandler)
		} else {
			err = k.receive(topics, offset, messageHandler)
		}
		if err != nil {
			log.Printf("consumer failed :%v", err)
		}
	}()
}

func (k *Kafka) receive(topics []string, offset int64, messageHandler MessageHandler) error {
	var err error
	k.ctx = context.Background()
	k.consumer, err = sarama.NewConsumer(k.brokers, k.config)
	if err != nil {
		return errors.New("new consumer failed: " + err.Error())
	}
	defer k.consumer.Close()
	for _, topic := range topics {
		go k.topic(topic, offset, messageHandler)
	}
	select {}
}

func (k *Kafka) receiveGroup(topics []string, offset int64, messageHandler MessageHandler) error {
	keepRunning := true
	k.config.Consumer.Offsets.Initial = offset
	ctx, cancel := context.WithCancel(context.Background())
	consumerGroup, err := sarama.NewConsumerGroup(k.brokers, k.groupId, k.config)
	if err != nil {
		log.Fatalf("new consumer group failed: %v", err)
		return err
	}
	defer consumerGroup.Close()
	consumer := Consumer{
		ready:          make(chan bool),
		messageHandler: messageHandler,
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

func (k *Kafka) topic(topic string, offset int64, messageHandler MessageHandler) {
	partitionList, err := k.consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("fail to start consumer partition,err:%v\n", err)
		return
	}
	for partition := range partitionList {
		k.partitionConsumer, err = k.consumer.ConsumePartition(topic, int32(partition), offset)
		if err != nil {
			log.Printf("fail to start consumer for partition %d,err:%v\n", partition, err)
			continue
		}
		go k.message(messageHandler)
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

func (k *Kafka) message(messageHandler MessageHandler) {
	defer k.partitionConsumer.AsyncClose()
	for message := range k.partitionConsumer.Messages() {
		messageHandler(message)
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
