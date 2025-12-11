package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/queue/pkg/factory"
	"github.com/yasin-wu/utils/queue/pkg/logger"
)

type rabbitMQ struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	logger      logger.Logger
	contentType string
	exchange    string
	direct      string
	forever     chan bool
	mandatory   bool
	immediate   bool
}

var _ factory.Queue = (*rabbitMQ)(nil)

func New(brokers []string, username, password string, config *config.MQConfig) (factory.Queue, error) {
	if len(brokers) == 0 {
		brokers = defaultBrokers
	}
	if config == nil {
		config = defaultConfig
	}
	if username == "" {
		username = defaultUsername
	}
	if password == "" {
		password = defaultPassword
	}
	uri := fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, brokers[0], config.VHost)
	rabbitmq := &rabbitMQ{
		contentType: defaultContentType,
		exchange:    defaultExchange,
		direct:      defaultDirect,
		forever:     make(chan bool),
		mandatory:   config.Mandatory,
		immediate:   config.Immediate,
		logger:      logger.NewDefaultLogger(),
	}
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	rabbitmq.conn = conn
	channel, err := rabbitmq.conn.Channel()
	if err != nil {
		return nil, err
	}

	rabbitmq.channel = channel
	return rabbitmq, nil
}
