package stdin

import (
	"bufio"
	"io"
	"os"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"sync"
)

var fact *factory

func init() {

	fact = &factory{
		ch: make(chan ql.Buffer),
	}

	ql.RegisterInput("stdin", fact)
}

// factory is the factory for the standard input
type factory struct {
	once sync.Once
	ch   chan ql.Buffer
}

// Build builds the input processor
func (f *factory) Build(_ io.Reader) (ql.InputProcess, error) {
	return &Process{}, nil
}

// Process is the standard input process
type Process struct {
}

// Start starts the standard input process
func (p *Process) Start(ctx context.Context, next chan<- ql.Buffer) error {

	log.Log(ctx).Debug("Starting input handler", "handler", "stdin")

	fact.once.Do(func() {
		go func() {
			bio := bufio.NewReader(os.Stdin)

			for {
				line, _, err := bio.ReadLine()
				if err != nil {
					break
				}
				fact.ch <- ql.Buffer{Data: line}
			}
		}()
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case str := <-fact.ch:
				next <- str
			}
		}
	}()

	return nil
}
