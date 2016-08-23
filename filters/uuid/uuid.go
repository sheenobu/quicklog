package uuid

import (
	"github.com/satori/go.uuid"
	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterFilter("uuid", &Handler{})
}

// Handler is the UUID handler
type Handler struct {
	FieldName string
}

// Handle is the quicklog handle method for processing a log line
func (u *Handler) Handle(ctx context.Context, prev <-chan ql.Line, next chan<- ql.Line, config map[string]interface{}) error {

	field := "uuid"
	if u.FieldName != "" {
		field = u.FieldName
	}

	ok := true

	fieldIface := config["field"]
	if fieldIface != nil {
		field, ok = fieldIface.(string)
		if !ok {
			log.Log(ctx).Warn("Could not parse UUID config, using field=uuid")
			field = "uuid"
		}
	}

	log.Log(ctx).Debug("Starting filter handler", "handler", "uuid", "field", field)

	go func() {
		for {
			select {
			case line := <-prev:
				line.Data[field] = uuid.NewV4().String()
				next <- line
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
