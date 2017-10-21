package goevents

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

// Event represents an event
type Event struct {
	Topic       string      `json:"topic"`
	ID          string      `json:"id"`
	EmittedTime time.Time   `json:"emitted_time"`
	Payload     interface{} `json:"payload"`
}

func (event *Event) MarshalJSON() ([]byte, error) {
	newEvent := *event
	rawPayload, err := json.Marshal(newEvent.Payload)
	if err != nil {
		return nil, err
	}
	newEvent.Payload = string(rawPayload)
	raw, err := json.Marshal(newEvent)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (event *Event) UnmarshalJSON(b []byte) error {
	var fields map[string]interface{}
	err := json.Unmarshal(b, &fields)
	if err != nil {
		return err
	}
	topic, ok := fields["topic"].(string)
	if !ok {
		return fmt.Errorf("Failed to decode field `topic`")
	}
	id, ok := fields["id"].(string)
	if !ok {
		return fmt.Errorf("Failed to decode field `id`")
	}
	timeStr, ok := fields["emitted_time"].(string)
	if !ok {
		return fmt.Errorf("Failed to decode field `emitted_time`")
	}
	time, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return fmt.Errorf("Failed to decode field `emitted_time`. %s", err)
	}
	event.Topic = topic
	event.ID = id
	event.EmittedTime = time
	event.Payload = fields["payload"]
	return nil
}

func (event *Event) Scan(payload EventPayload) error {
	rawStr := event.Payload.(string)
	raw := []byte(rawStr)
	return json.Unmarshal(raw, payload)
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
		ID:          uuid.NewV1().String(),
		EmittedTime: time.Now().Round(0),
		Payload:     payload,
	}
}
