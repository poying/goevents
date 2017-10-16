package nsq

import (
	gonsq "github.com/nsqio/go-nsq"
	"github.com/poying/goevents"
)

// NewBus creates an event bus
func NewBus(producer *gonsq.Producer, consumer *gonsq.Consumer) goevents.Bus {
	return goevents.NewBus(producer, &Consumer{consumer: consumer})
}

// Consumer represents a message consumer
type Consumer struct {
	consumer *gonsq.Consumer
}

// AddHandler adds a message handler to consumer
func (consumer *Consumer) AddHandler(handler goevents.MessageHandler) {
	consumer.consumer.AddHandler(gonsq.HandlerFunc(func(nsqMessage *gonsq.Message) error {
		defer nsqMessage.Finish()
		return handler.HandleMessage(&Message{message: nsqMessage})
	}))
}

// Message represents a message
type Message struct {
	message *gonsq.Message
}

// Decode the message
func (message *Message) Decode() ([]byte, error) {
	return message.message.Body, nil
}
