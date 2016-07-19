package stdout

import (
	"fmt"
	"os"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterOutput("stdout", &stdoutHandler{})
}

type stdoutHandler struct {
}

func (stdout *stdoutHandler) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {

	l := log.Log(ctx).New("handler", "stdout")

	l.Debug("Starting output handler", "handler", "stdout")

	go func() {
		for {
			select {
			case line := <-prev:
				if _, err := os.Stdout.Write([]byte(fmt.Sprintf("%s\n", line.Data["message"]))); err != nil {
					l.Error("Error writing to standard out", "error", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
