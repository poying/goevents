package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/poying/goevents"
	nsqevents "github.com/poying/goevents/nsq"
)

type Nested struct {
	Time time.Time
}

type Payload struct {
	Nested
	Message string
}

type EventHandler struct{}

func (handler *EventHandler) HandleEvent(event goevents.Event) error {
	payload := Payload{}
	err := event.Scan(&payload)
	if err != nil {
		return err
	}
	fmt.Println(payload.Time, payload.Message)
	return nil
}

func main() {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("topic", "channel", config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bus := nsqevents.NewBus(producer, consumer)

	bus.AddHandler(&EventHandler{})
	err = consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = bus.Publish("topic", Payload{
		Nested{Time: time.Now()},
		"Rocket Man",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	time.Sleep(10 * time.Second)
}
