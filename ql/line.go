package ql

// A Line is a distinct log line being processed by quicklog
type Line struct {
	Data map[string]string
}
