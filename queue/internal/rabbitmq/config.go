package rabbitmq

import (
	"github.com/yasin-wu/utils/queue/pkg/config"
	"github.com/yasin-wu/utils/util"
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

func (r *rabbitMQ) SetLogger(logger util.Logger) {
	if logger != nil {
		r.logger = logger
	}
}
