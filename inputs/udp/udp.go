package udp

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"net"
)

func init() {
	ql.RegisterInput("udp", ql.InputFactoryHandler(func(r io.Reader) (ql.InputProcess, error) {
		var config Config
		if err := json.NewDecoder(r).Decode(&config); err != nil {
			return nil, err
		}
		return &Process{Config: config}, nil
	}))
}

// Config defines the available configuration options for the udp input
type Config struct {
	Listen string `json:"listen"`
}

// Process is the process for the udp input handler
type Process struct {
	Config Config
}

// Start starts the UDP handler proces
func (p *Process) Start(ctx context.Context, next chan<- ql.Buffer) error {

	if p.Config.Listen == "" {
		return errors.New("No UDP listen configuration provided")
	}

	log.Log(ctx).Debug("Starting input handler", "handler", "udp", "listen", p.Config.Listen)

	ch := make(chan ql.Buffer)
	listenAddr, err := net.ResolveUDPAddr("udp", p.Config.Listen)
	if err != nil {
		return err
	}
	ln, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}

	go func() {

		buffer := make([]byte, 1024)

		for {

			size, addr, err := ln.ReadFromUDP(buffer)
			if err != nil {

				// if the context is done, assume the error
				// is that the listener has been closed.
				select {
				case <-ctx.Done():
					return
				default:
				}

				log.Log(ctx).Error("Error reading from connection", "error", err)

				continue
			}

			m := make(map[string]interface{})
			m["udp.source"] = addr.String()
			ch <- ql.Buffer{
				Data:     buffer[:size],
				Metadata: m,
			}
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
