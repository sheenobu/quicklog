package tcp

import (
	"bufio"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"net"
)

func init() {
	ql.RegisterInput("tcp", &tcpInput{})
}

type tcpInput struct {
}

func (tcp *tcpInput) Handle(ctx context.Context, next chan<- ql.Buffer, config map[string]interface{}) error {

	listen := config["listen"].(string)

	ch := make(chan ql.Buffer)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	log.Log(ctx).Debug("Starting input handler", "handler", "tcp")

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

					ch <- ql.CreateBuffer(line, m)
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
