package rabbitmq

import (
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/queue/pkg/logger"
)

var (
	defaultContentType = "application/json"
	defaultExchange    = "amq.direct"
	defaultDirect      = "direct"
	defaultBrokers     = []string{"localhost:5672"}
	defaultUsername    = "guest"
	defaultPassword    = "guest"
	defaultConfig      = &config.MQConfig{
		Mandatory: false,
		Immediate: false,
	}
)

func (r *rabbitMQ) SetLogger(logger logger.Logger) {
	if logger != nil {
		r.logger = logger
	}
}
