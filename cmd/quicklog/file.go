package main

import (
	"github.com/sheenobu/golibs/apps"
	"github.com/sheenobu/golibs/log"

	"github.com/sheenobu/quicklog/config"
	"github.com/sheenobu/quicklog/ql"

	"os"

	"golang.org/x/net/context"
)

func startFileQuicklog(ctx context.Context, app *apps.App) {

	// load config
	cfg, err := config.LoadFile(configFile)
	if err != nil {
		log.Log(ctx).Error("Error loading configuration", "error", err)
		os.Exit(255)
		return
	}

	// setup chain
	chain := ql.Chain{
		Input:        ql.GetInput(cfg.Input.Driver),
		InputConfig:  cfg.Input.Config,
		Parser:       ql.GetParser(cfg.Input.Parser),
		Output:       ql.GetOutput(cfg.Output.Driver),
		OutputConfig: cfg.Output.Config,
	}

	if len(cfg.Filters) >= 1 {
		chain.Filter = ql.GetFilter(cfg.Filters[0].Driver)
		chain.FilterConfig = cfg.Filters[0].Config
	}

	// execute chain
	app.SpawnSimple("chain", chain.Execute)

	return
}
