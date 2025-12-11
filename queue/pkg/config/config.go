package config

import (
	"github.com/Shopify/sarama"
	"github.com/nats-io/nats.go"
)

type KafkaConfig sarama.Config

type NatsConfig struct {
	Stream        string
	StreamEnabled bool
	Async         bool
	NCOpt         []nats.Option
	JSOpt         []nats.JSOpt
}

type MQConfig struct {
	ContentType string
	Exchange    string
	Direct      string
	VHost       string
	Mandatory   bool
	Immediate   bool
}
