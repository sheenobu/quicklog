package stdin

import (
	"bufio"
	"os"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"
)

func init() {
	ql.RegisterInput("stdin", &stdin{})
}

type stdin struct {
}

func (s *stdin) Handle(ctx context.Context, next chan<- ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting input handler", "handler", "stdin")

	ch := make(chan ql.Line)

	go func() {
		bio := bufio.NewReader(os.Stdin)
		for {
			line, _, err := bio.ReadLine()
			if err != nil {
				break
			}
			l := ql.Line{
				Data: make(map[string]string),
			}
			l.Data["message"] = string(line)
			ch <- l
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case str := <-ch:
				next <- str
			}
		}
	}()

	return nil
}
