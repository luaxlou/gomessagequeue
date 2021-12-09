package producerconsumer

import "github.com/luaxlou/gomessagequeue/mqengines"

type Producer struct {
	engine mqengines.MqEngine
}

func NewProducer(engine mqengines.MqEngine) *Producer {

	return &Producer{
		engine: engine,
	}
}

func (p *Producer) Produce(key string, data interface{}) error {

	return p.engine.Add(key, data)
}
