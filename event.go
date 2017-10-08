package goevents

import (
	"time"

	"github.com/satori/go.uuid"
)

// Event represents an event
type Event struct {
	Topic       string
	ID          []byte
	EmittedTime time.Time
	Payload     EventPayload
}

// EventPayload represents payload of event
type EventPayload interface{}

// EventHandler represents event handler
type EventHandler interface {
	HandleEvent(event Event) error
}

func newEvent(topic string, payload EventPayload) *Event {
	return &Event{
		Topic:       topic,
		ID:          uuid.NewV1().Bytes(),
		EmittedTime: time.Now().Round(0),
		Payload:     payload,
	}
}
