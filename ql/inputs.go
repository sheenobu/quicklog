package ql

import (
	"io"

	"golang.org/x/net/context"
)

var (
	inputs map[string]InputFactory
)

func init() {
	inputs = make(map[string]InputFactory)
}

// A Buffer is a series of characters and key-value data that get put on the pipeline
type Buffer struct {
	Data     []byte
	Metadata map[string]interface{}
}

// An InputProcess is a process that waits for input and sends it to the line channel
type InputProcess interface {
	Start(context.Context, chan<- Buffer) error
}

// InputProcessFunc is an adaptor for casting functions to InputProcess interfaces
type InputProcessFunc func(context.Context, chan<- Buffer) error

// Start starts the input processor by writing to the channel
func (i InputProcessFunc) Start(ctx context.Context, ch chan<- Buffer) error {
	return i(ctx, ch)
}

// InputFactory is the factory for building input processes given the config options
type InputFactory interface {
	Build(jsonConfig io.Reader) (InputProcess, error)
}

// InputFactoryHandler converts the given function to an InputFactory
type InputFactoryHandler func(io.Reader) (InputProcess, error)

// Build builds the input process given the JSON config reader
func (f InputFactoryHandler) Build(jsonConfig io.Reader) (InputProcess, error) {
	return f(jsonConfig)
}

// GetInput gets the input driver
func GetInput(driver string) InputFactory {
	return inputs[driver]
}

// RegisterInput registers the input handler using the driver name
func RegisterInput(driver string, factory InputFactory) {
	inputs[driver] = factory
}
