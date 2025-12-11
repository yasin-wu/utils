package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/queue/pkg/factory"
	"github.com/yasin-wu/utils/queue/pkg/logger"
)

type kafka struct {
	brokers           []string
	strategy          string
	forever           chan bool
	ctx               context.Context
	logger            logger.Logger
	config            *sarama.Config
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
}

var _ factory.Queue = (*kafka)(nil)

func New(brokers []string, _, _ string, config *config.KafkaConfig) (factory.Queue, error) {
	if config == nil {
		config = NewConfig()
	}
	if len(brokers) == 0 {
		brokers = defaultBrokers
	}
	kafka := &kafka{
		brokers: brokers,
		forever: make(chan bool),
		config:  (*sarama.Config)(config)}
	kafka.logger = logger.NewDefaultLogger()
	kafka.config.Consumer.Return.Errors = true
	kafka.config.Producer.Return.Successes = true
	switch kafka.strategy {
	case "sticky":
		kafka.config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		kafka.config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		kafka.config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		kafka.config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	}
	return kafka, nil
}
