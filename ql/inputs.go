package ql

import (
	"golang.org/x/net/context"
)

// An InputHandler is a handler that waits for input and sends it to the line channel
type InputHandler interface {
	Handle(chan<- Line, map[string]interface{}) (context.Context, error)
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
