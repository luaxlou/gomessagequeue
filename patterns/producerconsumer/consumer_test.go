package producerconsumer

import (
	"log"
	"testing"
)

func TestNewConsumer(t *testing.T) {

	engine := getTestEngine()

	c := NewConsumer(engine)

	c.Consume("key", 1000, func(contents []string) error {

		log.Println(contents)

		return nil
	})

}
