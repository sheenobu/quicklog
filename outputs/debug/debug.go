package debug

import (
	"fmt"
	"os"

	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterOutput("debug", &debugHandler{})
}

type debugHandler struct {
}

func (d *debugHandler) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {

	printFields := true

	pfInput := config["print-fields"]
	if pfInput != nil {
		switch k := pfInput.(type) {
		case bool:
			printFields = k
		case string:
			printFields = k == "true"
		default:
			log.Log(ctx).Warn("Could not parse print-fields variable, falling back to true")
		}
	}

	log.Log(ctx).Debug("Starting output handler", "handler", "debug")

	go func() {
		for {
			select {
			case line := <-prev:
				os.Stdout.Write([]byte(fmt.Sprintf("Time: [%v]\n", line.Timestamp)))
				if line.Data["message"] == nil {
					line.Data["message"] = ""
				}

				os.Stdout.Write([]byte("Message: " + line.Data["message"].(string) + "\n"))

				if !printFields {
					continue
				}

				for key, val := range line.Data {
					if key != "message" {
						os.Stdout.Write([]byte(fmt.Sprintf("\t%s=%s\n", key, val)))
					}
				}

				os.Stdout.Write([]byte("\n"))
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
