package ql

import (
	"golang.org/x/net/context"
)

// A FilterHandler is a handler that processes, edits, and filters incoming data
type FilterHandler interface {
	Handle(context.Context, <-chan Line, chan<- Line, map[string]interface{}) error
}

var (
	filters map[string]FilterHandler
)

func init() {
	filters = make(map[string]FilterHandler)
}

// GetFilter returns the filter handler
func GetFilter(driver string) FilterHandler {
	return filters[driver]
}

// RegisterFilter registers the filter handler under the driver name
func RegisterFilter(driver string, handler FilterHandler) {
	filters[driver] = handler
}
