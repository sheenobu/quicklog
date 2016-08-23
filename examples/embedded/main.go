package main

import (
	"github.com/sheenobu/quicklog/filters/uuid"
	"github.com/sheenobu/quicklog/inputs/stdin"
	"github.com/sheenobu/quicklog/outputs/debug"
	"github.com/sheenobu/quicklog/parsers/plain"

	"golang.org/x/net/context"

	"github.com/sheenobu/quicklog/ql"
)

func main() {

	chain := ql.Chain{
		Input: &stdin.Process{},
		//Output: &stdout.Process{},
		Output: &debug.Handler{PrintFields: debug.NullableBool{NotNull: false, Value: true}},
		Filter: &uuid.Handler{FieldName: "uuid"},
		Parser: &plain.Parser{},
	}

	ctx := context.Background()
	chain.Execute(ctx)

}
