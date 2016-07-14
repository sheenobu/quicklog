package ql

import (
	"golang.org/x/net/context"
)

// A Buffer is a series of characters and key-value data that get put on the pipeline
type Buffer struct {
	Data     []byte
	Metadata map[string]interface{}
}

// An InputHandler is a handler that waits for input and sends it to the line channel
type InputHandler interface {
	Handle(context.Context, chan<- Buffer, map[string]interface{}) error
}

var (
	inputs map[string]InputHandler
)

func init() {
	inputs = make(map[string]InputHandler)
}

// GetInput gets the input driver
func GetInput(driver string) InputHandler {
	return inputs[driver]
}

// RegisterInput registers the input handler using the driver name
func RegisterInput(driver string, handler InputHandler) {
	inputs[driver] = handler
}
