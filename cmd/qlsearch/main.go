package main

// simple index search for bleve

import (
	"fmt"
	"github.com/blevesearch/bleve"
)

func main() {
	// open a new index
	index, err := bleve.Open("example.bleve")

	fmt.Printf("error: %s\n", err)

	// search for some text
	query := bleve.NewMatchAllQuery()
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	fmt.Printf("error: %s\n", err)

	fmt.Printf("%s\n", searchResults)

}
