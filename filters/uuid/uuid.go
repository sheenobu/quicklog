package uuid

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"
)

func init() {
	ql.RegisterFilter("uuid", &uuidHandler{})
}

type uuidHandler struct {
}

func (u *uuidHandler) Handle(ctx context.Context, prev chan ql.Line, next chan ql.Line, config map[string]interface{}) error {

	field := "uuid"
	ok := true

	fieldIface := config["field"]
	if fieldIface != nil {
		field, ok = fieldIface.(string)
		if !ok {
			//TODO: Warn
			field = "uuid"
		}
	}

	go func() {
		for {
			select {
			case line := <-prev:
				line.Data[field] = uuid.New()
				next <- line
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
