package goevents

import (
	"bytes"
	"encoding/gob"
	"sync"
)

// Bus represents event bus
type Bus struct {
	lock     *sync.RWMutex
	producer Producer
	consumer Consumer
	handlers []EventHandler
}

// NewBus creates an event bus
func NewBus(producer Producer, consumer Consumer) *Bus {
	bus := &Bus{
		lock:     &sync.RWMutex{},
		producer: producer,
		consumer: consumer,
	}
	consumer.AddHandler(bus)
	return bus
}

// Publish publishs an event
func (bus *Bus) Publish(topic string, payload EventPayload) error {
	event := newEvent(topic, payload)
	var body bytes.Buffer
	enc := gob.NewEncoder(&body)
	err := enc.Encode(event)
	if err != nil {
		return err
	}
	return bus.producer.Publish(topic, body.Bytes())
}

// HandleMessage decodes messages which are received from message queue and pass them to event handlers
func (bus *Bus) HandleMessage(message Message) error {
	var event Event
	body, err := message.Decode()
	gobDecoder := gob.NewDecoder(bytes.NewReader(body))
	err = gobDecoder.Decode(&event)
	if err != nil {
		return err
	}
	bus.lock.RLock()
	defer bus.lock.RUnlock()
	for _, handler := range bus.handlers {
		err = handler.HandleEvent(event)
		if err != nil {
			return err
		}
	}
	return nil

}

// AddHandler adds an event handler
func (bus *Bus) AddHandler(handler EventHandler) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.handlers = append(bus.handlers, handler)
}
