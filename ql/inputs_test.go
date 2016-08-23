package ql

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestInput(t *testing.T) {
	fact := GetInput("dummy")
	if fact != nil {
		t.Error("Dummy input should have been nil")
	}

	RegisterInput("dummy", InputFactoryHandler(func(r io.Reader) (InputProcess, error) {
		return nil, errors.New("Error building dummy input")
	}))

	fact = GetInput("dummy")
	if fact == nil {
		t.Error("Dummy input should not have been nil")
	}
	proc, err := fact.Build(strings.NewReader(""))
	if proc != nil {
		t.Error("Dummy process should have been nil")
	}
	if err.Error() != "Error building dummy input" {
		t.Error("Expected error to be 'Error building dummy input'")
	}
}
