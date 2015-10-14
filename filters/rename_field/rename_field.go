package uuid

import (
	"fmt"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterFilter("rename_field", &rename{})
}

type rename struct {
}

func (*rename) Handle(ctx context.Context, prev <-chan ql.Line, next chan<- ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting filter handler", "handler", "rename_field")

	src, exists := config["source"].(string)
	if !exists {
		return fmt.Errorf("Error getting source config in rename_field handler")
	}

	dst, exists := config["dest"].(string)
	if !exists {
		return fmt.Errorf("Error getting dest config in rename_field handler")
	}

	cp, exists := config["copy"].(bool)
	if !exists {
		cp = false
	}

	go func() {
		for {
			select {
			case line := <-prev:
				line.Data[dst] = line.Data[src]
				if !cp { // empty out source field
					line.Data[src] = nil
				}
				next <- line
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
