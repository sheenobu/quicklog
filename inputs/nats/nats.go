package nats

import (
	"github.com/nats-io/nats"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"fmt"
)

// read each line from a nats queue

func init() {
	ql.RegisterInput("nats", &natsIn{})
}

type natsIn struct{}

func (in *natsIn) Handle(ctx context.Context, next chan<- ql.Buffer, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting input handler", "handler", "nats")

	url, ok := config["url"].(string)
	if !ok || url == "" {
		log.Log(ctx).Error("Could not create nats input, no url defined")
		return fmt.Errorf("Could not create nats input, no url defined")
	}

	opts := nats.DefaultOptions
	opts.Url = url

	servers, ok := config["servers"].([]string)
	if ok {
		opts.Servers = servers
	}

	publish, ok := config["subscribe"].(string)
	if !ok || publish == "" {
		log.Log(ctx).Error("Could not create nats input, no publish defined")
		return fmt.Errorf("Could not create nats input, no publish defined")
	}

	encoding, ok := config["encoding"].(string)
	if !ok || encoding == "" {
		encoding = nats.JSON_ENCODER
	}

	nc, err := opts.Connect()
	if err != nil {
		log.Log(ctx).Error("Error connecting to nats url", "url", url, "error", err)
		return err
	}

	c, err := nats.NewEncodedConn(nc, encoding)
	if err != nil {
		log.Log(ctx).Error("Error creating nats connection", "error", err)
		return err
	}

	recvCh := make(chan map[string]interface{})
	c.BindRecvChan(publish, recvCh)

	go func() {
		defer c.Close()

		for {
			select {
			case line := <-recvCh:
				msg := ""
				if line["message"] != nil {
					msg = line["message"].(string)
				}

				delete(line, "message")
				next <- ql.CreateBuffer([]byte(msg), line)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
