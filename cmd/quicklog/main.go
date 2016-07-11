package main

import (
	_ "github.com/sheenobu/quicklog/filters"
	_ "github.com/sheenobu/quicklog/inputs"
	_ "github.com/sheenobu/quicklog/outputs"
	_ "github.com/sheenobu/quicklog/parsers"

	"github.com/sheenobu/golibs/managed"
	"github.com/sheenobu/quicklog/log"

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
