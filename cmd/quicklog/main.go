package main

import (
	"fmt"

	"github.com/sheenobu/quicklog/config"

	_ "github.com/sheenobu/quicklog/filters"
	_ "github.com/sheenobu/quicklog/inputs"
	_ "github.com/sheenobu/quicklog/outputs"

	"github.com/sheenobu/quicklog/ql"
)

func main() {

	cfg, err := config.LoadFile("quicklog.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	chain := ql.Chain{
		Input:       ql.GetInput(cfg.Input.Driver),
		InputConfig: cfg.Input.Config,

		Filter:       ql.GetFilter(cfg.Filters[0].Driver),
		FilterConfig: cfg.Filters[0].Config,

		Output:       ql.GetOutput(cfg.Output.Driver),
		OutputConfig: cfg.Output.Config,
	}

	chain.Execute()
}
