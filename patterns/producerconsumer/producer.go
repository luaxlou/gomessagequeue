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

//考虑到数据的多样性，使用string兼容所有结构，如果是对象类型，可以序列化后传输
func (p *Producer) Produce(key string, content string, maxLen int64) error {

	return p.engine.Add(key, content, maxLen)
}
