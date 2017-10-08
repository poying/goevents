package local_test

import (
	"testing"
	"time"

	"github.com/bouk/monkey"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/poying/goevents"
	"github.com/poying/goevents/local"
	mocks "github.com/poying/goevents/mocks"
)

type Payload struct {
	Message string
}

func TestBus(t *testing.T) {
	id := uuid.NewV1()
	now := time.Now().Round(0)

	timePatcher := monkey.Patch(time.Now, func() time.Time { return now })
	uuidPatcher := monkey.Patch(uuid.NewV1, func() uuid.UUID { return id })
	defer timePatcher.Unpatch()
	defer uuidPatcher.Unpatch()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bus := local.NewBus()
	bus.RegisterPayloadType(Payload{})
	handler := mocks.NewMockEventHandler(ctrl)
	bus.AddHandler(handler)
	event := goevents.Event{
		Topic:       "hello",
		ID:          id.Bytes(),
		Payload:     Payload{Message: "Hello"},
		EmittedTime: now,
	}
	handler.EXPECT().HandleEvent(event).Times(1)
	err := bus.Publish("hello", event.Payload)
	assert.Nil(t, err)
}
