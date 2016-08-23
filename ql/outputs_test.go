package ql

import (
	"errors"
	"testing"

	"golang.org/x/net/context"
)

func TestOutput(t *testing.T) {
	handler := GetOutput("dummy")
	if handler != nil {
		t.Errorf("Dummy handler should be nil")
	}

	RegisterOutput("dummy", OutputHandlerFunc(func(ctx context.Context, ch <-chan Line, cfg map[string]interface{}) error {
		return errors.New("Error running output handler")
	}))

	handler = GetOutput("dummy")
	if handler == nil {
		t.Error("Dummy handler should not be nil")
		return
	}

	err := handler.Handle(nil, nil, nil)
	if err.Error() != "Error running output handler" {
		t.Error("Expected error to be 'Error running output handler'")
	}

}
