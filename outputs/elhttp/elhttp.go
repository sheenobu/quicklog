package nats

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"fmt"
)

// publish the json to elasticsearch

func init() {
	ql.RegisterOutput("elasticsearch-http", &elHTTP{})
}

type elHTTP struct{}

func (out *elHTTP) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting output handler", "handler", "elasticsearch-http")

	url, ok := config["url"].(string)
	if !ok || url == "" {
		log.Log(ctx).Error("Could not create elasticsearch-http output, no url defined")
		return fmt.Errorf("Could not create elasticsearch-http output, no url defined")
	}

	index, ok := config["index"].(string)
	if !ok || index == "" {
		index = "quicklog"
	}

	_type, ok := config["type"].(string)
	if !ok || _type == "" {
		_type = "entry"
	}

	go func() {
		for {
			select {
			case line := <-prev:
				data := line.Data
				data["timestamp"] = line.Timestamp

				//TODO: eventually provide functions that can determine index, _type from line data

				destURL := fmt.Sprintf("%s/%s/%s", url, index, _type)

				r, w := io.Pipe()

				go func() {
					defer w.Close()
					err := json.NewEncoder(w).Encode(&data)
					if err != nil {
						log.Log(ctx).Error("error converting line to json", "error", err)
					}
				}()

				resp, err := http.Post(destURL, "application/json", r)
				if err != nil {
					log.Log(ctx).Error("error sending to elasticsearch", "error", err)
					continue
				}
				defer resp.Body.Close()

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
