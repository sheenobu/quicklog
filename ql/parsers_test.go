package ql

import (
	"errors"
	"testing"
)

func TestParser(t *testing.T) {
	handler := GetParser("dummy")
	if handler != nil {
		t.Errorf("Dummy parser should be nil")
	}

	RegisterParser("dummy", ParserFunc(func(buffer []byte, line *Line, cfg map[string]interface{}) error {
		return errors.New("Error running parser")
	}))

	handler = GetParser("dummy")
	if handler == nil {
		t.Error("Dummy parser should not be nil")
		return
	}

	err := handler.Parse(nil, nil, nil)
	if err.Error() != "Error running parser" {
		t.Error("Expected error to be 'Error running parser'")
	}

}
