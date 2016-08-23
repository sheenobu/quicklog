package ql

import (
	"golang.org/x/net/context"
)

var (
	filters map[string]FilterHandler
)

func init() {
	filters = make(map[string]FilterHandler)
}

// A FilterHandler is a handler that processes, edits, and filters incoming data
type FilterHandler interface {
	Handle(context.Context, <-chan Line, chan<- Line, map[string]interface{}) error
}

// FilterHandlerFunc is the adaptor for converting a function to a FilterHandler instance
type FilterHandlerFunc func(context.Context, <-chan Line, chan<- Line, map[string]interface{}) error

// Handle runs the filter by reading from in, processing the line, and writing it to out
func (fh FilterHandlerFunc) Handle(ctx context.Context, in <-chan Line, out chan<- Line, cfg map[string]interface{}) error {
	return fh(ctx, in, out, cfg)
}

// GetFilter returns the filter handler
func GetFilter(driver string) FilterHandler {
	return filters[driver]
}

// RegisterFilter registers the filter handler under the driver name
func RegisterFilter(driver string, handler FilterHandler) {
	filters[driver] = handler
}
