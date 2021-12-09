package producerconsumer

import (
	"github.com/luaxlou/gomessagequeue/mqengines"
	"github.com/luaxlou/gomessagequeue/mqengines/redismqengine"
	"strconv"
	"testing"
	"time"
)

func getTestEngine() mqengines.MqEngine {
	engine, err := redismqengine.New("127.0.0.1:6379", "", 0)

	if err != nil {
		panic(err.Error())

	}

	return engine
}

func TestProducer_Produce(t *testing.T) {

	engine := getTestEngine()

	p := NewProducer(engine)

	go func() {

		for i := 0; i < 100000; i++ {
			p.Produce("key", strconv.Itoa(i))

			time.Sleep(time.Second * 1)
		}

	}()

	select {

	}

}
