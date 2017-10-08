package local

import (
	"github.com/poying/goevents"
)

// NewBus creates an event bus
func NewBus() *goevents.Bus {
	consumer := &Consumer{}
	producer := &Producer{consumer: consumer}
	return goevents.NewBus(producer, consumer)
}

// Producer represents a message producer
type Producer struct {
	consumer *Consumer
}

// Publish publishs an event
func (producer *Producer) Publish(topic string, body []byte) error {
	producer.consumer.handleMessage(&Message{body: body})
	return nil
}

// Consumer represents a message consumer
type Consumer struct {
	handlers []goevents.MessageHandler
}

// AddHandler adds a message handler to consumer
func (consumer *Consumer) AddHandler(handler goevents.MessageHandler) {
	consumer.handlers = append(consumer.handlers, handler)
}

func (consumer *Consumer) handleMessage(message goevents.Message) {
	for _, handler := range consumer.handlers {
		handler.HandleMessage(message)
	}
}

// Message represents a message
type Message struct {
	body []byte
}

// Decode the message
func (message *Message) Decode() ([]byte, error) {
	return message.body, nil
}
