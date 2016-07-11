package log

import (
	"testing"

	"golang.org/x/net/context"
)

func TestEmptyContext(t *testing.T) {

	l := Log(context.Background())

	if l == nil {
		t.Errorf("Expected log to be not empty")
	}

	l.Debug("hello")

}

func TestSimpleContext(t *testing.T) {

	ctx := context.Background()
	ctx = NewContext(ctx, "source", "test")

	l := Log(ctx)

	if l == nil {
		t.Errorf("Expected log to be not empty")
	}

	l.Debug("hello")

}
