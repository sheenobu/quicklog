package hostname

import (
	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
	"os"
)

func init() {
	ql.RegisterFilter("hostname", &hostnameHandler{})
}

type hostnameHandler struct {
}

func (*hostnameHandler) Handle(ctx context.Context, prev <-chan ql.Line, next chan<- ql.Line, config map[string]interface{}) error {

	field := "hostname"
	ok := true

	fieldIface := config["field"]
	if fieldIface != nil {
		field, ok = fieldIface.(string)
		if !ok {
			log.Log(ctx).Warn("Could not parse hostname config, using field=hostname")
			field = "hostname"
		}
	}

	log.Log(ctx).Debug("Starting filter handler", "handler", "hostname", "field", field)

	hostname, _ := os.Hostname()

	go func() {
		for {
			select {
			case line := <-prev:
				line.Data[field] = hostname
				next <- line
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
