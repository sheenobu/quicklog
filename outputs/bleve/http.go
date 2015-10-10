package bleve

import (
	bleveHttp "github.com/blevesearch/bleve/http"
	"net/http"
)

func (out *bleveOutput) startHttpServer(listen string) {
	bleveHttp.RegisterIndexName("example", out.index)
	searchHandler := bleveHttp.NewSearchHandler("example")

	http.Handle("/search", searchHandler)

	assets := http.StripPrefix("/", http.FileServer(http.Dir("static/")))
	http.Handle("/", assets)
	http.ListenAndServe(listen, nil)
}
