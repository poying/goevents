package goevents

import "encoding/gob"

// RegisterPayloadType regiters payload type to payload encoder/decoder
func RegisterPayloadType(payloadType interface{}) {
	gob.Register(payloadType)
}
