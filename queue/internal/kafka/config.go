package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/util"
)

var defaultBrokers = []string{"localhost:9092"}

func NewConfig() *config.KafkaConfig {
	conf := (*config.KafkaConfig)(sarama.NewConfig())
	conf.Version = sarama.MaxVersion
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	return conf
}

func (k *kafka) SetLogger(logger util.Logger) {
	if logger != nil {
		k.logger = logger
	}
}
