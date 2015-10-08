package tcp

import (
	"bufio"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"net"
	"time"
)

func init() {
	ql.RegisterInput("tcp", &tcpInput{})
}

type tcpInput struct {
}

func (tcp *tcpInput) Handle(ctx context.Context, next chan<- ql.Line, config map[string]interface{}) error {

	listen := config["listen"].(string)

	ch := make(chan ql.Line)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	log.Log(ctx).Debug("Starting input handler", "handler", "tcp")

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Log(ctx).Error("Error accepting connection", "error", err)
			}
			go func(conn net.Conn) {
				bio := bufio.NewReader(conn)
				defer conn.Close()
				for {
					line, _, err := bio.ReadLine()
					if err != nil {
						break
					}
					l := ql.Line{
						Data: make(map[string]string),
					}

					l.Timestamp = time.Now() //TODO: read timestamp from incoming data
					l.Data["message"] = string(line)
					l.Data["tcp.source"] = conn.RemoteAddr().String()
					ch <- l
				}
			}(conn)
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
