package kafka

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type Kafka struct {
	brokers           []string
	groupId           string
	strategy          string
	ctx               context.Context
	config            *sarama.Config
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
	messageHandler    MessageHandler
}

type MessageHandler func(message *sarama.ConsumerMessage)

type Option func(kafka *Kafka)

type Config sarama.Config

func New(brokers []string, config *Config, options ...Option) *Kafka {
	if config == nil {
		config = NewConfig()
	}
	kafka := &Kafka{brokers: brokers, config: (*sarama.Config)(config)}
	for _, f := range options {
		f(kafka)
	}
	if kafka.messageHandler == nil {
		kafka.messageHandler = printMsg
	}
	return kafka
}

func NewConfig() *Config {
	return (*Config)(sarama.NewConfig())
}

func WithGroupId(groupId string) Option {
	return func(kafka *Kafka) {
		kafka.groupId = groupId
	}
}

func WithStrategy(strategy string) Option {
	return func(kafka *Kafka) {
		kafka.strategy = strategy
	}
}

func WithMessageHandler(messageHandler MessageHandler) Option {
	return func(kafka *Kafka) {
		kafka.messageHandler = messageHandler
	}
}

func (k *Kafka) SetGroupId(groupId string) {
	k.groupId = groupId
}

func (k *Kafka) SetStrategy(strategy string) {
	k.strategy = strategy
}

func (k *Kafka) SetMessageHandler(messageHandler MessageHandler) {
	k.messageHandler = messageHandler
}
