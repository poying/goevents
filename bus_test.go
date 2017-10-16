package goevents_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/poying/goevents"
	mocks "github.com/poying/goevents/mocks"
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
