package goevents

// Message represents message
type Message interface {
	Decode() ([]byte, error)
}

// Producer represents message producer
type Producer interface {
	Publish(topic string, body []byte) error
}

// Consumer represents message consumer
type Consumer interface {
	AddHandler(handler MessageHandler)
}

// MessageHandler represents message handler
type MessageHandler interface {
	HandleMessage(message Message) error
}
