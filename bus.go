package goevents

import (
	"encoding/json"
	"sync"
)

// Bus represents event bus
type Bus interface {
	Publish(topic string, payload EventPayload) error
	AddHandler(handler EventHandler)
}

type bus struct {
	lock     *sync.RWMutex
	producer Producer
	consumer Consumer
	handlers []EventHandler
}

// NewBus creates an event bus
func NewBus(producer Producer, consumer Consumer) Bus {
	b := &bus{
		lock:     &sync.RWMutex{},
		producer: producer,
		consumer: consumer,
	}
	consumer.AddHandler(&funcHandler{handler: b.handleMessage})
	return b
}

// Publish publishs an event
func (b *bus) Publish(topic string, payload EventPayload) error {
	event := newEvent(topic, payload)
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return b.producer.Publish(topic, body)
}

func (b *bus) handleMessage(message Message) error {
	var event Event
	body, err := message.Decode()
	err = json.Unmarshal(body, &event)
	if err != nil {
		return err
	}
	b.lock.RLock()
	defer b.lock.RUnlock()
	for _, handler := range b.handlers {
		err = handler.HandleEvent(event)
		if err != nil {
			return err
		}
	}
	return nil

}

// AddHandler adds an event handler
func (b *bus) AddHandler(handler EventHandler) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.handlers = append(b.handlers, handler)
}

type funcHandler struct {
	handler func(message Message) error
}

func (handler *funcHandler) HandleMessage(message Message) error {
	return handler.handler(message)
}
