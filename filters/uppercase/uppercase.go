package uuid

import (
	"strings"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterFilter("uppercase", &uppercase{})
}

type uppercase struct {
}

func (u *uppercase) Handle(ctx context.Context, prev <-chan ql.Line, next chan<- ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting filter handler", "handler", "uppercase")

	go func() {
		for {
			select {
			case line := <-prev:
				line.Data["message"] = strings.ToUpper(line.Data["message"].(string))
				next <- line
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
