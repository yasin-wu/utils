package queue

import (
	"errors"
	"strings"

	"github.com/yasin-wu/utils/queue/internal/kafka"
	"github.com/yasin-wu/utils/queue/internal/nats"
	"github.com/yasin-wu/utils/queue/internal/rabbitmq"
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/queue/pkg/factory"
	"github.com/yasin-wu/utils/util"
)

type Config struct {
	Brokers        []string
	Driver         string
	Username       string
	Password       string
	Logger         util.Logger
	KafkaConfig    *config.KafkaConfig
	RabbitMQConfig *config.MQConfig
	NatsConfig     *config.NatsConfig
}

type Option func(client *Client)

type Client struct {
	factory.Queue
}

func New(config *Config, options ...Option) (*Client, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	var (
		err   error
		queue factory.Queue
	)
	switch strings.ToLower(config.Driver) {
	case "kafka":
		queue, err = kafka.New(config.Brokers, config.Username, config.Password, config.KafkaConfig)
	case "rabbitmq":
		queue, err = rabbitmq.New(config.Brokers, config.Username, config.Password, config.RabbitMQConfig)
	case "nats":
		queue, err = nats.New(config.Brokers, config.Username, config.Password, config.NatsConfig)
	default:
		return nil, errors.New("not supported driver")
	}
	if err != nil {
		return nil, err
	}
	queue.SetLogger(config.Logger)
	client := &Client{Queue: queue}
	for _, f := range options {
		f(client)
	}
	return client, nil
}
