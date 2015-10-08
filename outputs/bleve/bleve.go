package bleve

import (
	"github.com/blevesearch/bleve"
	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"code.google.com/p/go-uuid/uuid"
)

func init() {
	ql.RegisterOutput("bleve", &bleveOutput{})
}

type bleveOutput struct{}

func (out *bleveOutput) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting output handler", "handler", "bleve")

	index, err := bleve.Open("example.bleve")
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New("example.bleve", mapping)
		if err != nil {
			return err
		}
	}

	go func() {
		for {
			select {
			case line := <-prev:
				err = index.Index(uuid.New(), line.Data)
				if err != nil {
					log.Log(ctx).Error("Error indexing line", "error", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
