package main

import (
	"github.com/sheenobu/quicklog/config"

	_ "github.com/sheenobu/quicklog/filters"
	_ "github.com/sheenobu/quicklog/inputs"
	_ "github.com/sheenobu/quicklog/outputs"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/ql"

	"golang.org/x/net/context"

	"flag"
	"os"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "filename", "quicklog.json", "Filename for the configuration")
}

func main() {

	flag.Parse()

	ctx := context.Background()
	ctx = log.NewContext(ctx)
	log.Log(ctx).Info("Starting quicklog", "configfile", configFile)

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
		Output:       ql.GetOutput(cfg.Output.Driver),
		OutputConfig: cfg.Output.Config,
	}

	if len(cfg.Filters) >= 1 {
		chain.Filter = ql.GetFilter(cfg.Filters[0].Driver)
		chain.FilterConfig = cfg.Filters[0].Config
	}

	// execute chain
	chain.Execute(ctx)
}
