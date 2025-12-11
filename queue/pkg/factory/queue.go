package factory

import (
	"github.com/yasin-wu/utils/queue/pkg/consumer"
	"github.com/yasin-wu/utils/queue/pkg/logger"
	"github.com/yasin-wu/utils/queue/pkg/message"
)

type Queue interface {
	//Topics() ([]string, error)
	Publish(messages ...*message.Message) error
	Subscribe(group string, consumers ...*consumer.Consumer)
	SetLogger(logger logger.Logger)
	Stop()
}
