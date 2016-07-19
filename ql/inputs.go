package ql

import (
	"io"

	"golang.org/x/net/context"
)

// A Buffer is a series of characters and key-value data that get put on the pipeline
type Buffer struct {
	Data     []byte
	Metadata map[string]interface{}
}

// An InputProcess is a process that waits for input and sends it to the line channel
type InputProcess interface {
	Start(context.Context, chan<- Buffer) error
}

// InputFactory is the factory for building input processes given the config options
type InputFactory interface {
	Build(jsonConfig io.Reader) (InputProcess, error)
}

// InputFactoryHandler converts the given function to an InputFactory
func InputFactoryHandler(fn func(io.Reader) (InputProcess, error)) InputFactory {
	return &factfn{fn}
}

type factfn struct {
	fn func(io.Reader) (InputProcess, error)
}

func (f *factfn) Build(jsonConfig io.Reader) (InputProcess, error) {
	return f.fn(jsonConfig)
}

var (
	inputs map[string]InputFactory
)

func init() {
	inputs = make(map[string]InputFactory)
}

// GetInput gets the input driver
func GetInput(driver string) InputFactory {
	return inputs[driver]
}

// RegisterInput registers the input handler using the driver name
func RegisterInput(driver string, factory InputFactory) {
	inputs[driver] = factory
}
