package tcp

import (
	"bufio"

	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"net"
)

func init() {
	ql.RegisterInput("tcp", &tcpInput{})
}

type tcpInput struct {
}

func (tcp *tcpInput) Handle(next chan<- ql.Line, config map[string]interface{}) (context.Context, error) {

	ctx, _ := context.WithCancel(context.Background())

	listen := config["listen"].(string)

	ch := make(chan ql.Line)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return ctx, err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				//TODO: log error
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

	return ctx, nil
}
