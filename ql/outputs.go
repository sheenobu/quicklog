package ql

import (
	"golang.org/x/net/context"
)

// OutputHandler is the interface for defining output plugins
type OutputHandler interface {
	Handle(context.Context, <-chan Line, map[string]interface{}) error
}

var (
	outputs map[string]OutputHandler
)

func init() {
	outputs = make(map[string]OutputHandler)
}

// GetOutput returns the output handler for the given driver
func GetOutput(driver string) OutputHandler {
	return outputs[driver]
}

// RegisterOutput registers an output handler given the driver name
func RegisterOutput(driver string, handler OutputHandler) {
	outputs[driver] = handler
}
