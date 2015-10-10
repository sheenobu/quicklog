package bleve

import (
	"github.com/blevesearch/bleve"
	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"
	"golang.org/x/net/context"

	"code.google.com/p/go-uuid/uuid"

	"sync"
)

func init() {
	ql.RegisterOutput("bleve", &bleveOutput{})
}

type bleveOutput struct {
	index bleve.Index
	once  sync.Once
}

func (out *bleveOutput) Handle(ctx context.Context, prev <-chan ql.Line, config map[string]interface{}) error {

	log.Log(ctx).Debug("Starting output handler", "handler", "bleve")

	out.once.Do(func() {
		var err error
		out.index, err = bleve.Open("example.bleve")
		if err != nil {
			mapping := bleve.NewIndexMapping()
			out.index, err = bleve.New("example.bleve", mapping)
			if err != nil {
				return
			}
		}

		listen := ":8080"

		if config["http.listen"] != nil {
			listen = config["http.listen"].(string)
		}

		go out.startHttpServer(listen)
	})

	go func() {
		for {
			select {
			case line := <-prev:
				line.Data["timestamp"] = line.Timestamp
				err := out.index.Index(uuid.New(), line.Data)
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
