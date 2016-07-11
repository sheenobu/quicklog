// Package log implements logging attached to golang.org/x/net/context
package log

import (
	"golang.org/x/net/context"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

type logKeyType int

var logKey logKeyType

// SetLog sets the log on the context
func SetLog(ctx context.Context, l log15.Logger) context.Context {
	return context.WithValue(ctx, logKey, l)
}

// NewContext creates a context with the log populated
func NewContext(ctx context.Context, args ...interface{}) context.Context {
	l := log15.New(args...)
	return context.WithValue(ctx, logKey, l)
}

// FromContext gets the logger out of the context
func FromContext(ctx context.Context) (log15.Logger, bool) {
	l, ok := ctx.Value(logKey).(log15.Logger)
	return l, ok
}

// Log gets the log from the context
func Log(ctx context.Context) log15.Logger {
	l, ok := FromContext(ctx)
	if !ok {
		l = log15.New()
	}
	return l
}
