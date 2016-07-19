package nats

import (
	"encoding/json"
	"io"

	"github.com/nats-io/nats"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"fmt"
)

func init() {
	ql.RegisterInput("nats", ql.InputFactoryHandler(func(r io.Reader) (ql.InputProcess, error) {
		var config Config
		if err := json.NewDecoder(r).Decode(&config); err != nil {
			return nil, err
		}
		return &Process{Config: config}, nil
	}))
}

// Config defines the available configuration options for the nats input
type Config struct {
	URL      string   `json:"url"`
	Servers  []string `json:"servers"`
	Publish  string   `json:"subscribe"`
	Encoding string   `json:"encoding"`
}

// Process defines the nats input object
type Process struct {
	Config Config
}

// Start reads from the nats queue and pushes onto the given next channel
func (i *Process) Start(ctx context.Context, next chan<- ql.Buffer) error {

	log.Log(ctx).Debug("Starting input handler", "handler", "nats")

	if i.Config.URL == "" {
		log.Log(ctx).Error("Could not create nats input, no url defined")
		return fmt.Errorf("Could not create nats input, no url defined")
	}

	opts := nats.DefaultOptions
	opts.Url = i.Config.URL
	opts.Servers = i.Config.Servers

	if i.Config.Publish == "" {
		log.Log(ctx).Error("Could not create nats input, no publish defined")
		return fmt.Errorf("Could not create nats input, no publish defined")
	}

	if i.Config.Encoding == "" {
		i.Config.Encoding = nats.JSON_ENCODER
	}

	nc, err := opts.Connect()
	if err != nil {
		log.Log(ctx).Error("Error connecting to nats url", "url", i.Config.URL, "error", err)
		return err
	}

	c, err := nats.NewEncodedConn(nc, i.Config.Encoding)
	if err != nil {
		log.Log(ctx).Error("Error creating nats connection", "error", err)
		return err
	}

	recvCh := make(chan ql.Line)
	sub, err := c.BindRecvChan(i.Config.Publish, recvCh)
	if err != nil {
		log.Log(ctx).Error("Error listening on nats receive channel", "error", err)
		return err
	}

	go func() {
		defer sub.Unsubscribe()
		defer c.Close()

		for {
			select {
			case line := <-recvCh:
				var msg string
				if m, ok := line.Data["message"]; ok {
					msg = m.(string)
				}

				delete(line.Data, "message")
				next <- ql.Buffer{
					Data:     []byte(msg),
					Metadata: line.Data,
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
