package nats

import (
	"errors"
	"fmt"

	"github.com/yasin-wu/utils/util"

	natsmodel "github.com/nats-io/nats.go"
)

var defaultStream = "queue"

func (n *nats) SetLogger(logger util.Logger) {
	if logger != nil {
		n.logger = logger
	}
}

func (n *nats) addStream(name string) error {
	streamConfig := &natsmodel.StreamConfig{
		Name:     name,
		Subjects: []string{fmt.Sprintf("%s.>", name)},
	}
	stream, err := n.jetStream.StreamInfo(name)
	if err != nil && !errors.Is(err, natsmodel.ErrStreamNotFound) {
		n.logger.Errorf("stream info failed : %v", err)
		return err
	}
	if stream != nil {
		if _, err = n.jetStream.UpdateStream(streamConfig); err != nil {
			n.logger.Errorf("add stream failed : %v", err)
			return err
		}
	}
	if stream == nil || errors.Is(err, natsmodel.ErrStreamNotFound) {
		if _, err = n.jetStream.AddStream(streamConfig); err != nil {
			n.logger.Errorf("add stream failed : %v", err)
			return err
		}
	}
	return nil
}
