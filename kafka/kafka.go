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
	strategy          Strategy
	ctx               context.Context
	config            *sarama.Config
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
	messageHandler    MessageHandler
}

type Strategy string

const (
	Sticky_Strategy     Strategy = "sticky"
	Roundrobin_Strategy Strategy = "roundrobin"
	Range_Strategy      Strategy = "range"
)

type MessageHandler func(message *sarama.ConsumerMessage)

type Option func(kafka *Kafka)

type Config sarama.Config

/**
 * @author: yasin
 * @date: 2022/2/9 11:48
 * @params: brokers []string, config *Config, options ...Option
 * @return: *Kafka
 * @description: new kafka client
 */
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

/**
 * @author: yasin
 * @date: 2022/2/9 11:50
 * @return: *Config
 * @description: new kafka config
 */
func NewConfig() *Config {
	return (*Config)(sarama.NewConfig())
}

/**
 * @author: yasin
 * @date: 2022/2/9 11:48
 * @params: groupId string
 * @return: Option
 * @description: new kafka client with groupId
 */
func WithGroupId(groupId string) Option {
	return func(kafka *Kafka) {
		kafka.groupId = groupId
	}
}

/**
 * @author: yasin
 * @date: 2022/2/9 11:49
 * @params: strategy Strategy
 * @return: Option
 * @description: new kafka client with strategy
 */
func WithStrategy(strategy Strategy) Option {
	return func(kafka *Kafka) {
		kafka.strategy = strategy
	}
}

/**
 * @author: yasin
 * @date: 2022/2/9 11:49
 * @params: messageHandler MessageHandler
 * @return: Option
 * @description: new kafka client with messageHandler
 */
func WithMessageHandler(messageHandler MessageHandler) Option {
	return func(kafka *Kafka) {
		kafka.messageHandler = messageHandler
	}
}

/**
 * @author: yasin
 * @date: 2022/2/9 11:49
 * @params: groupId string
 * @description: set groupId
 */
func (k *Kafka) SetGroupId(groupId string) {
	k.groupId = groupId
}

/**
 * @author: yasin
 * @date: 2022/2/9 11:50
 * @params: strategy Strategy
 * @description: set strategy
 */
func (k *Kafka) SetStrategy(strategy Strategy) {
	k.strategy = strategy
}

/**
 * @author: yasin
 * @date: 2022/2/9 11:50
 * @params: messageHandler MessageHandler
 * @description: set messageHandler
 */
func (k *Kafka) SetMessageHandler(messageHandler MessageHandler) {
	k.messageHandler = messageHandler
}
