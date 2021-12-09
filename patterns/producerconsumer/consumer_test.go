package producerconsumer

import (
	"log"
	"testing"
)

func TestConsumer_Consume(t *testing.T) {

	engine := getTestEngine()

	c := NewConsumer(engine)

	c.Consume("key", 1000, func(contents []string) error {

		log.Println(contents)

		return nil
	})

}

func TestConsumer_Consume_Black(t *testing.T) {

	engine := getTestEngine()

	c := NewConsumer(engine)


	go func() {
		c.ConsumeBlock("key", 1000, func(contents []string) error {

			log.Println(contents)

			return nil
		})

	}()


	go func() {
		c.ConsumeBlock("key", 1000, func(contents []string) error {

			log.Println(contents)

			return nil
		})

	}()

	select {

	}

}
