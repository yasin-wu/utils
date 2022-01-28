package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func printMsg(msg *sarama.ConsumerMessage) {
	fmt.Printf("this is consumer message, Topic:%s Partition:%d Offset:%d Key:%s Value:%s\n",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
}
