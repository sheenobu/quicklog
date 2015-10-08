package stdout

import (
	"fmt"
	"os"

	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterOutput("debug", &debugHandler{})
}

type debugHandler struct {
}

func (d *debugHandler) Handle(ctx context.Context, prev chan ql.Line, config map[string]interface{}) error {

	printFields := true

	pfInput := config["print-fields"]
	if pfInput != nil {
		switch k := pfInput.(type) {
		case bool:
			printFields = k
		case string:
			printFields = k == "true"
		default:
			//TODO: log warning
		}
	}

	go func() {
		for {
			select {
			case line := <-prev:
				os.Stdout.Write([]byte("Message: " + line.Data["message"]))

				if !printFields {
					os.Stdout.Write([]byte("\n"))
					continue
				} else {
					os.Stdout.Write([]byte(" | "))
				}

				for key, val := range line.Data {
					if key != "message" {
						os.Stdout.Write([]byte(fmt.Sprintf("%s=%s ", key, val)))
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
