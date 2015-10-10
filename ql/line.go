package ql

import (
	"time"
)

// A Line is a distinct log line being processed by quicklog
type Line struct {
	Timestamp time.Time
	Data      map[string]interface{}
}
