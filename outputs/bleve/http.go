package bleve

import (
	bleveHttp "github.com/blevesearch/bleve/http"
	"net/http"

	"github.com/sheenobu/golibs/log"
	"golang.org/x/net/context"
)

func (out *bleveOutput) startHttpServer(ctx context.Context, cfg map[string]interface{}) {

	listen := ":8080"

	if cfg["listen"] != nil {
		listen = cfg["listen"].(string)
	}

	root := ""

	if cfg["webroot"] != nil {
		root = cfg["webroot"].(string)
	}

	bleveHttp.RegisterIndexName("example", out.index)
	searchHandler := bleveHttp.NewSearchHandler("example")

	http.Handle("/search", searchHandler)

	log.Log(ctx).Info("Starting bleve HTTP search server", "listen", listen)

	if root != "" {
		log.Log(ctx).Debug("Serving up static fileserver", "root", root)
		assets := http.StripPrefix("/", http.FileServer(http.Dir(root)))
		http.Handle("/", assets)
	}

	http.ListenAndServe(listen, nil)
}
