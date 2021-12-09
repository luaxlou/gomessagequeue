package producerconsumer

import (
	"github.com/luaxlou/gomessagequeue/mqengines"
)

type Consumer struct {
	engine mqengines.MqEngine
}

func NewConsumer(engine mqengines.MqEngine) *Consumer {

	return &Consumer{
		engine: engine,
	}
}

func (c *Consumer) Consume(key string, count int64, onRead func(contents []string) error) {

	c.engine.Read(key, count, onRead)
}
