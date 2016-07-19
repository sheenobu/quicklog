package managed

import (
	"fmt"
	"io"
	"strconv"

	"time"

	"golang.org/x/net/context"
)

// A Component is our fundamental building block
type Component func(ctx context.Context)

// Simple defines a simple component which runs and blocks on context completion
func Simple(name string, c Component) *Process {
	return &Process{
		Name:      name,
		Component: c,
		data: map[string]string{
			"Type":   "simple",
			"Name":   name,
			"Status": "stopped",
		},
		Writer: func(process *Process, w io.Writer) error {
			t, _ := process.Get("Type")
			n, _ := process.Get("Name")
			s, _ := process.Get("Status")
			_, err := w.Write([]byte(fmt.Sprintf("%s: %s [%s]\n", t, n, s)))
			return err
		},
	}
}

// Timer defines a timer component which runs on the given interval
func Timer(name string, interval time.Duration, runImmediately bool, f func(context.Context)) *Process {

	process := &Process{
		Name: name,
		data: map[string]string{
			"Type":     "timer",
			"Name":     name,
			"Status":   "stopped",
			"Interval": interval.String(),
		},
		Writer: func(process *Process, w io.Writer) error {
			t, _ := process.Get("Type")
			n, _ := process.Get("Name")
			i, _ := process.Get("Interval")
			l, _ := process.Get("Lastrun")
			s, _ := process.Get("Status")
			r, _ := process.Get("Runcount")

			_, err := w.Write([]byte(fmt.Sprintf("%s: %s runcount=%s interval=%s lastrun=%s [%s]\n", t, n, r, i, l, s)))
			return err
		},
	}

	process.Component = func(ctx context.Context) {
		runCount := 0
		if runImmediately {
			runCount++
			process.Set("Runcount", strconv.Itoa(runCount))
			process.Set("Lastrun", time.Now().String())
			f(ctx)
		}
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-time.After(interval):
				runCount++
				process.Set("Lastrun", t.String())
				process.Set("Runcount", strconv.Itoa(runCount))
				f(ctx)
			}
		}
	}

	return process
}
