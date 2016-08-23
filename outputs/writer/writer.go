package writer

import (
	"errors"
	"fmt"
	"io"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

// Process is the writer output process
type Process struct {
	Name string
	W    io.Writer
}

// Handle is the quicklog output handler
func (p *Process) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {

	if p.W == nil {
		return errors.New("No writer provided to the output handler")
	}

	if p.Name == "" {
		p.Name = "writer"
	}

	l := log.Log(ctx).New("handler", p.Name)

	l.Debug("Starting output handler")

	go func() {
		for {
			select {
			case line := <-prev:
				if _, err := p.W.Write([]byte(fmt.Sprintf("%s\n", line.Data["message"]))); err != nil {
					l.Error("Error writing", "error", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
