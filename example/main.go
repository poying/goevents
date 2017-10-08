package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/poying/goevents"
	nsqevents "github.com/poying/goevents/nsq"
)

type Payload struct {
	Message string
}

type EventHandler struct{}

func (handler *EventHandler) HandleEvent(event goevents.Event) error {
	switch payload := event.Payload.(type) {
	case Payload:
		fmt.Println(payload.Message)
	}
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

	bus.RegisterPayloadType(Payload{})
	bus.AddHandler(&EventHandler{})
	err = consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = bus.Publish("topic", Payload{Message: "Rocket Man"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	time.Sleep(10 * time.Second)
}
