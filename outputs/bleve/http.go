package bleve

import (
	"encoding/json"
	"net/http"

	"github.com/blevesearch/bleve"
	"github.com/sheenobu/quicklog/log"
	"golang.org/x/net/context"
)

func (out *bleveOutput) startHTTPServer(ctx context.Context, cfg map[string]interface{}) {

	listen := ":8080"

	if cfg["listen"] != nil {
		listen = cfg["listen"].(string)
	}

	root := ""

	if cfg["webroot"] != nil {
		root = cfg["webroot"].(string)
	}

	http.Handle("/search", http.HandlerFunc(out.searchHandler))

	log.Log(ctx).Info("Starting bleve HTTP search server", "listen", listen)

	if root != "" {
		log.Log(ctx).Debug("Serving up static fileserver", "root", root)
		assets := http.StripPrefix("/", http.FileServer(http.Dir(root)))
		http.Handle("/", assets)
	}

	http.ListenAndServe(listen, nil)
}

// SearchRequest is the JSON payload for the search request
type SearchRequest struct {
	Size  int    `json:"size"`
	Query string `json:"query"`
	From  int    `json:"from"`
}

func (out *bleveOutput) searchHandler(w http.ResponseWriter, r *http.Request) {

	var req SearchRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err) //TODO: handle error properly
	}

	if req.Query != "" && req.Query != "*" {

		sr := bleve.SearchRequest{
			Query:   bleve.NewQueryStringQuery(req.Query),
			From:    req.From,
			Size:    req.Size,
			Fields:  []string{"*"},
			Explain: false,
		}

		searchResults, err := out.index.Search(&sr)
		if err != nil {
			panic(err)
		}

		if err := json.NewEncoder(w).Encode(&searchResults); err != nil {
			panic(err)
		}
	} else {

		sr := bleve.SearchRequest{
			Query:   bleve.NewMatchAllQuery(),
			From:    req.From,
			Size:    req.Size,
			Fields:  []string{"*"},
			Explain: false,
		}

		searchResults, err := out.index.Search(&sr)
		if err != nil {
			panic(err)
		}

		if err := json.NewEncoder(w).Encode(&searchResults); err != nil {
			panic(err)
		}
	}
}
