package goevents_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/poying/goevents"
	mocks "github.com/poying/goevents/mocks"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	consumer := mocks.NewMockConsumer(ctrl)
	producer := mocks.NewMockProducer(ctrl)
	consumer.EXPECT().AddHandler(gomock.Any()).Times(1)
	producer.EXPECT().Publish("hello", gomock.Any()).Times(1)
	bus := goevents.NewBus(producer, consumer)
	err := bus.Publish("hello", 123)
	assert.Nil(t, err)
}

func TestHandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	consumer := mocks.NewMockConsumer(ctrl)
	producer := mocks.NewMockProducer(ctrl)
	consumer.EXPECT().AddHandler(gomock.Any()).Times(1)
	bus := goevents.NewBus(producer, consumer)

	event := goevents.Event{
		Topic:       "hello",
		ID:          uuid.NewV1().Bytes(),
		Payload:     123,
		EmittedTime: time.Now().Round(0),
	}
	var encodedEvent bytes.Buffer
	enc := gob.NewEncoder(&encodedEvent)
	err := enc.Encode(event)
	assert.Nil(t, err)

	message := mocks.NewMockMessage(ctrl)
	message.EXPECT().Decode().Return(encodedEvent.Bytes(), nil).Times(1)

	fmt.Println(event)
	handler := mocks.NewMockEventHandler(ctrl)
	handler.EXPECT().HandleEvent(event).Times(1)
	bus.AddHandler(handler)

	bus.HandleMessage(message)
}
