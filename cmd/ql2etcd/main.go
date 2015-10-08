package main

// ql2etcd pushes a quicklog json config file to etcd

import (
	"github.com/sheenobu/golibs/log"

	"github.com/sheenobu/quicklog/config"

	"os"

	"flag"
	"golang.org/x/net/context"
)

var (
	configFile    string
	etcdEndpoints string
	instanceName  string
)

func init() {
	flag.StringVar(&configFile, "input", "", "Filename for the configuration")

	flag.StringVar(&etcdEndpoints, "etcdEndpoints", "", "Servers for etcd, comma separated")
	flag.StringVar(&instanceName, "instanceName", "", "Instance name used for etcd prefix")

}

func main() {

	flag.Parse()

	// Setup context
	ctx := context.Background()
	ctx = log.NewContext(ctx)

	// load config
	cfg, err := config.LoadFile(configFile)
	if err != nil {
		log.Log(ctx).Error("Error loading configuration", "error", err)
		os.Exit(255)
		return
	}

	if err := syncToEtcd(ctx, cfg); err != nil {
		log.Log(ctx).Error("Error writing configuration into etcd", "error", err)
		os.Exit(255)
		return
	}

	os.Exit(0)
}
