package ql

import (
	"errors"
	"testing"

	"golang.org/x/net/context"
)

func TestFilter(t *testing.T) {
	handler := GetFilter("dummy")
	if handler != nil {
		t.Errorf("Dummy filter should be nil")
	}

	RegisterFilter("dummy", FilterHandlerFunc(func(ctx context.Context, ch <-chan Line, out chan<- Line, cfg map[string]interface{}) error {
		return errors.New("Error running filter")
	}))

	handler = GetFilter("dummy")
	if handler == nil {
		t.Error("Dummy filter should not be nil")
		return
	}

	err := handler.Handle(nil, nil, nil, nil)
	if err.Error() != "Error running filter" {
		t.Error("Expected error to be 'Error running filter'")
	}

}
