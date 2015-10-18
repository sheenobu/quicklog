package udp

import (
	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"net"
)

func init() {
	ql.RegisterInput("udp", &udpInput{})
}

type udpInput struct {
}

func (udp *udpInput) Handle(ctx context.Context, next chan<- ql.Buffer, config map[string]interface{}) error {

	listen := config["listen"].(string)
	log.Log(ctx).Debug("Starting input handler", "handler", "udp", "listen", listen)

	ch := make(chan ql.Buffer)
	listenAddr, err := net.ResolveUDPAddr("udp", listen)
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
			ch <- ql.CreateBuffer(buffer[:size], m)
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
