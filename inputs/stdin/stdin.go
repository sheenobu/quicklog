package stdin

import (
	"bufio"
	"os"
	"time"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"sync"
)

func init() {
	ql.RegisterInput("stdin", &stdin{
		ch: make(chan ql.Line),
	})
}

type stdin struct {
	once sync.Once
	ch   chan ql.Line
}

func (s *stdin) Handle(ctx context.Context, next chan<- ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting input handler", "handler", "stdin")

	s.once.Do(func() {
		go func() {
			bio := bufio.NewReader(os.Stdin)

			for {

				line, _, err := bio.ReadLine()
				if err != nil {
					break
				}
				l := ql.Line{
					Data:      make(map[string]string),
					Timestamp: time.Now(),
				}
				l.Data["message"] = string(line)
				s.ch <- l
			}
		}()
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case str := <-s.ch:
				next <- str
			}
		}
	}()

	return nil
}
