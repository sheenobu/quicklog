package tcp

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"net"
)

func init() {
	ql.RegisterInput("tcp", ql.InputFactoryHandler(func(r io.Reader) (ql.InputProcess, error) {
		var config Config
		if err := json.NewDecoder(r).Decode(&config); err != nil {
			return nil, err
		}
		return &Process{Config: config}, nil
	}))
}

// Config defines the available configuration options for the tcp input
type Config struct {
	Listen string `json:"listen"`
}

// Process defines the tcp process listener
type Process struct {
	Config Config
}

// Start starts the tcp listener and waits for messages
func (p *Process) Start(ctx context.Context, next chan<- ql.Buffer) error {

	if p.Config.Listen == "" {
		return errors.New("No TCP listen option provided")
	}

	log.Log(ctx).Debug("Starting input handler", "handler", "tcp", "listen", p.Config.Listen)

	ch := make(chan ql.Buffer)
	ln, err := net.Listen("tcp", p.Config.Listen)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {

				// if the context is done, assume the error
				// is that the listener has been closed.
				select {
				case <-ctx.Done():
					return
				default:
				}

				log.Log(ctx).Error("Error accepting connection", "error", err)

				continue
			}
			go func(conn net.Conn) {
				bio := bufio.NewReader(conn)
				for {
					line, _, err := bio.ReadLine()
					if err != nil {
						break
					}

					m := make(map[string]interface{})
					m["tcp.source"] = conn.RemoteAddr().String()

					ch <- ql.Buffer{
						Data:     line,
						Metadata: m,
					}
				}
			}(conn)
		}
	}()

	go func() {
		defer ln.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case buffer := <-ch:
				next <- buffer
			}
		}
	}()

	return nil
}
