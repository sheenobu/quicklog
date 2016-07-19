package main

import (
	"github.com/sheenobu/golibs/managed"
	"github.com/sheenobu/quicklog/log"

	"github.com/sheenobu/quicklog/config"

	"os"

	"golang.org/x/net/context"
)

func startFileQuicklog(ctx context.Context, system *managed.System) {

	log.Log(ctx).Info("Loading config from file", "file", configFile)

	// load config
	cfg, err := config.LoadFile(configFile)
	if err != nil {
		log.Log(ctx).Error("Error loading configuration", "error", err)
		os.Exit(255)
		return
	}

	// setup chain
	chain := fromConfig(ctx, cfg)

	// execute chain
	system.Add(managed.Simple("chain", chain.Execute))

	return
}
