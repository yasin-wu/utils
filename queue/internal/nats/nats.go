package nats

import (
	"fmt"
	"strings"

	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/queue/pkg/factory"
	"github.com/yasin-wu/utils/queue/pkg/logger"

	natsmodel "github.com/nats-io/nats.go"
)

type nats struct {
	async         bool
	stream        string
	streamEnabled bool
	forever       chan bool
	logger        logger.Logger
	conn          *natsmodel.Conn
	jetStream     natsmodel.JetStreamContext
}

var _ factory.Queue = (*nats)(nil)

func New(brokers []string, username, password string, conf *config.NatsConfig) (factory.Queue, error) {
	url := natsmodel.DefaultURL
	if len(brokers) > 0 {
		url = strings.Join(brokers, ",")
	}
	if conf == nil {
		conf = &config.NatsConfig{}
	}
	if username != "" && password != "" {
		conf.NCOpt = append(conf.NCOpt, natsmodel.UserInfo(username, password))
	}
	conf.NCOpt = append(conf.NCOpt, natsmodel.RetryOnFailedConnect(true))
	conf.JSOpt = append(conf.JSOpt, natsmodel.PublishAsyncMaxPending(256))
	nats := &nats{
		async:         conf.Async,
		stream:        defaultStream,
		streamEnabled: conf.StreamEnabled,
		forever:       make(chan bool),
		logger:        logger.NewDefaultLogger(),
	}
	if conf.Stream != "" {
		nats.stream = conf.Stream
	}
	conn, err := natsmodel.Connect(url, conf.NCOpt...)
	if err != nil {
		return nats, fmt.Errorf("connect nats failed: %v", err)
	}
	nats.conn = conn
	return nats.initStream(conf)
}

func (n *nats) initStream(config *config.NatsConfig) (factory.Queue, error) {
	if n.streamEnabled {
		var err error
		if n.jetStream, err = n.conn.JetStream(config.JSOpt...); err != nil {
			return n, fmt.Errorf("connect jetstream failed: %v", err)
		}
		if err = n.addStream(n.stream); err != nil {
			return n, fmt.Errorf("add jetstream failed: %v", err)
		}
	}
	return n, nil
}
