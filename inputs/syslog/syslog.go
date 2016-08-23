package syslog

import (
	"encoding/json"
	"errors"
	"io"
	"net"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterInput("syslog", ql.InputFactoryHandler(func(r io.Reader) (ql.InputProcess, error) {
		var config Config
		if err := json.NewDecoder(r).Decode(&config); err != nil {
			return nil, err
		}
		return &Process{Config: config}, nil
	}))
}

// Config defines the available configuration options for the syslog input
type Config struct {
	Listen string `json:"listen"`
}

// Process is the syslog listener process
type Process struct {
	Config Config
}

const timestamp = "Mmm dd hh:mm:ss"

// Start starts the syslog handler proces
func (p *Process) Start(ctx context.Context, next chan<- ql.Buffer) error {

	if p.Config.Listen == "" {
		return errors.New("No syslog listen configuration provided")
	}

	log.Log(ctx).Debug("Starting input handler", "handler", "syslog", "listen", p.Config.Listen)

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

			priStart := 0
			priEnd := 0

			for _, i := range buffer {
				priEnd++
				if i == '>' {
					break
				}
			}

			hostStart := priEnd + len(timestamp) + 1
			hostEnd := priEnd + len(timestamp) + 1

			for _, i := range buffer[hostStart:] {
				if i == ' ' {
					break
				}
				hostEnd++

			}

			tagStart := hostEnd + 1
			tagEnd := hostEnd + 1
			for _, i := range buffer[tagStart:] {
				if i == ':' || i == '[' || i == ' ' {
					break
				}
				tagEnd++
			}

			m := make(map[string]interface{})
			m["syslog.tag"] = string(buffer[tagStart:tagEnd])
			m["syslog.hostname"] = string(buffer[hostStart:hostEnd])
			m["syslog.timestamp"] = string(buffer[priEnd : priEnd+len(timestamp)])
			m["syslog.pri"] = string(buffer[priStart+1 : priEnd-1])
			m["udp.source"] = addr.String()

			data := buffer[tagEnd+2 : size]

			ch <- ql.Buffer{
				Data:     data,
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
