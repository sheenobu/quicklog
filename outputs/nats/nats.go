package nats

import (
	"github.com/nats-io/nats"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"fmt"
)

// publish each line to a nats queue

func init() {
	ql.RegisterOutput("nats", &natsOut{})
}

type natsOut struct{}

func (out *natsOut) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting output handler", "handler", "nats")

	url, ok := config["url"].(string)
	if !ok || url == "" {
		log.Log(ctx).Error("Could not create nats output, no url defined")
		return fmt.Errorf("Could not create nats output, no url defined")
	}

	opts := nats.DefaultOptions
	opts.Url = url

	servers, ok := config["servers"].([]string)
	if ok {
		opts.Servers = servers
	}

	publish, ok := config["publish"].(string)
	if !ok || publish == "" {
		log.Log(ctx).Error("Could not create nats output, no publish defined")
		return fmt.Errorf("Could not create nats output, no publish defined")
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

	go func() {
		defer c.Close()

		for {
			select {
			case line := <-prev:
				err = c.Publish(publish, line)
				if err != nil {
					log.Log(ctx).Error("Error publishing to nats connection", "error", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
