package kafka

import (
	"context"
	"os"

	"github.com/Shopify/sarama"
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/queue/pkg/factory"
	"github.com/yasin-wu/utils/util"
)

type kafka struct {
	brokers []string
	//strategy          string
	forever           chan bool
	ctx               context.Context
	logger            util.Logger
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
	kafka.logger = util.NewDefaultLogger()
	if os.Getenv("KAFKA_DEBUG") == "true" {
		sarama.Logger = &kafkaLogger{Logger: kafka.logger}
	}
	kafka.config.Consumer.Return.Errors = true
	kafka.config.Producer.Return.Successes = true
	if len(kafka.config.Consumer.Group.Rebalance.GroupStrategies) == 0 {
		kafka.config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	}
	return kafka, nil
}

type kafkaLogger struct {
	util.Logger
}

var _ sarama.StdLogger = (*kafkaLogger)(nil)

func (l *kafkaLogger) Print(v ...interface{}) {
	l.Infof("kafka log:%v", v...)
}
func (l *kafkaLogger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}
func (l *kafkaLogger) Println(v ...interface{}) {
	l.Infof("kafka log:%v", v...)
}
