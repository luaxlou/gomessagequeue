package producerconsumer

import (
	"log"
	"testing"
)

func TestNewConsumer(t *testing.T) {

	engine := getTestEngine()

	c := NewConsumer(engine)

	c.Consume("key", 1000, func(data []interface{}) error {

		log.Println(data)

		return nil
	})

}
