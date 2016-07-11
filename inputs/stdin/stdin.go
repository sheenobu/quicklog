package stdin

import (
	"bufio"
	"os"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"sync"
)

func init() {
	ql.RegisterInput("stdin", &stdin{
		ch: make(chan ql.Buffer),
	})
}

type stdin struct {
	once sync.Once
	ch   chan ql.Buffer
}

func (s *stdin) Handle(ctx context.Context, next chan<- ql.Buffer, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting input handler", "handler", "stdin")

	s.once.Do(func() {
		go func() {
			bio := bufio.NewReader(os.Stdin)

			for {
				line, _, err := bio.ReadLine()
				if err != nil {
					break
				}
				s.ch <- ql.CreateBuffer(line, make(map[string]interface{}))
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
