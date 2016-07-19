package main

import (
	_ "github.com/sheenobu/quicklog/filters/uuid"
	_ "github.com/sheenobu/quicklog/inputs/stdin"
	_ "github.com/sheenobu/quicklog/outputs/stdout"
	_ "github.com/sheenobu/quicklog/parsers/plain"

	"golang.org/x/net/context"

	"github.com/sheenobu/quicklog/ql"
)

func main() {

	chain := ql.Chain{
		Input:  ql.GetInput("stdin"),
		Output: ql.GetOutput("stdout"),
		Filter: ql.GetFilter("uuid"),
		Parser: ql.GetParser("plain"),
	}

	ctx := context.Background()
	chain.Execute(ctx)

}
