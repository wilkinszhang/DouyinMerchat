package mq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var OrderProducer rocketmq.Producer

func InitProducer() error {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
	)
	if err != nil {
		return err
	}
	OrderProducer = p
	return p.Start()
}

func ShutdownProducer() {
	OrderProducer.Shutdown()
}
