package main

import (
	_ "github.com/sheenobu/quicklog/filters"
	_ "github.com/sheenobu/quicklog/inputs"
	_ "github.com/sheenobu/quicklog/outputs"

	"github.com/sheenobu/golibs/apps"
	"github.com/sheenobu/golibs/log"

	"golang.org/x/net/context"

	"flag"
	"os"
	"os/signal"
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

	// Setup app
	app := apps.NewApp("quicklog")
	app.StartWithContext(ctx)

	// register signal listeners
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		log.Log(ctx).Info("Got interrupt signal, stopping quicklog", "signal", s)
		app.Stop()
	}()

	switch {
	case etcdEndpoints != "" && instanceName != "":
		startEtcdQuicklog(ctx, app)
	case configFile != "":
		startFileQuicklog(ctx, app)
	}

	app.Wait()
}
