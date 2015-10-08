package stdout

import (
	"fmt"
	"os"

	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterOutput("stdout", &stdoutHandler{})
}

type stdoutHandler struct {
}

func (stdout *stdoutHandler) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {
	go func() {
		for {
			select {
			case line := <-prev:
				os.Stdout.Write([]byte(fmt.Sprintf("%s\n", line.Data["message"])))
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
