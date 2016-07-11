package main

import (
	_ "github.com/sheenobu/quicklog/filters"
	_ "github.com/sheenobu/quicklog/inputs"
	_ "github.com/sheenobu/quicklog/outputs"
	_ "github.com/sheenobu/quicklog/parsers"

	"github.com/sheenobu/quicklog/config"
	"github.com/sheenobu/quicklog/log"
	"github.com/sheenobu/quicklog/ql"

	"github.com/sheenobu/golibs/managed"

	"golang.org/x/net/context"

	"flag"
	"os"
)

var configFile string
var etcdEndpoints string
var instanceName string

func init() {

	flag.StringVar(&configFile, "filename", "quicklog.json", "Filename for the configuration")

	flag.StringVar(&etcdEndpoints, "etcdEndpoints", "", "Servers for etcd, comma separated")
	flag.StringVar(&instanceName, "instanceName", "", "Instance name used for etcd prefix")
}

func main() {

	flag.Parse()

	// Setup context
	ctx := context.Background()
	ctx = log.NewContext(ctx)
	log.Log(ctx).Info("Starting quicklog")

	// Setup system
	system := managed.NewSystem("quicklog")
	system.StartWithContext(ctx)

	// Register signal listeners
	system.RegisterForStop(os.Interrupt, os.Kill)

	switch {
	case etcdEndpoints != "" && instanceName != "":
		startEtcdQuicklog(ctx, system)
	case configFile != "":
		startFileQuicklog(ctx, system)
	}

	system.Wait()
}

func fromConfig(cfg *config.Config) *ql.Chain {
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

	return &chain
}
