package main

import (
	"encoding/json"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"github.com/sheenobu/golibs/apps"
	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/config"
	"github.com/sheenobu/quicklog/ql"

	"strconv"
	"strings"
	"time"
)

func startEtcdQuicklog(mainCtx context.Context, app *apps.App) {

	log.Log(mainCtx).Debug("Connecting to endpoint", "endpoints", strings.Split(etcdEndpoints, ","))
	etcdCfg := client.Config{
		Endpoints: strings.Split(etcdEndpoints, ","),
	}

	c, err := client.New(etcdCfg)
	if err != nil {
		log.Log(mainCtx).Crit("Error connecting to etcd server", "error", err)
		return
	}
	root := "/quicklog/" + instanceName

	log.Log(mainCtx).Debug("Listening for etcd keys", "key", root)

	kapi := client.NewKeysAPI(c)

	chainApp := apps.NewApp("app-chain-" + instanceName)
	app.SpawnApp(chainApp)

	var cfg config.Config

	err = syncFromEtcd(mainCtx, root, kapi, &cfg)
	if err != nil {
		log.Log(mainCtx).Error("Error syncing from etcd", "error", err)
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

	chainApp.SpawnSimple("chain-sub-"+instanceName, chain.Execute)

	app.SpawnSimple("etcd", func(ctx context.Context) {

		w := kapi.Watcher(root+"/reload", &client.WatcherOptions{
			Recursive: false,
		})

		for {
			resp, err := w.Next(ctx)

			if err != nil {
				if err == context.DeadlineExceeded {
					continue
				} else if err == context.Canceled {
					return
				} else if cerr, ok := err.(*client.ClusterError); ok {
					for _, e := range cerr.Errors {
						if e != context.Canceled {
							log.Log(ctx).Error("Error getting next etcd watch event", "parentError", err, "error", e)
						}
					}
				} else {
					log.Log(ctx).Error("Error getting next etcd watch event", "error", err)
				}

				return
			}
			if resp == nil {
				return
			}

			switch resp.Action {
			case "get":
				// do nothing
			default:
				log.Log(ctx).Info("Got update on quicklog config", "etcd.action", resp.Action)

				var newCfg config.Config

				err = syncFromEtcd(ctx, root, kapi, &newCfg)
				if err != nil {
					log.Log(ctx).Error("Error syncing from etcd", "error", err)
				} else {

					chainApp.Stop()
					<-time.After(1 * time.Second)
					chainApp = apps.NewApp("app-chain-" + instanceName)
					app.SpawnApp(chainApp)

					// setup chain
					chain = ql.Chain{
						Input:        ql.GetInput(newCfg.Input.Driver),
						InputConfig:  newCfg.Input.Config,
						Output:       ql.GetOutput(newCfg.Output.Driver),
						OutputConfig: newCfg.Output.Config,
					}

					if len(newCfg.Filters) >= 1 {
						chain.Filter = ql.GetFilter(newCfg.Filters[0].Driver)
						chain.FilterConfig = newCfg.Filters[0].Config
					}

					chainApp.SpawnSimple("chain-sub-"+instanceName, chain.Execute)
				}
				//TODO: (re)load config
			}
		}
	})
}

func syncFromEtcd(ctx context.Context, root string, cl client.KeysAPI, cfg *config.Config) error {
	inputDriverResponse, err := cl.Get(ctx, root+"/input/driver", nil)
	if err != nil {
		return err
	}

	inputDriverCfg, err := cl.Get(ctx, root+"/input/config", nil)
	if err != nil {
		return err
	}

	outputDriverResponse, err := cl.Get(ctx, root+"/output/driver", nil)
	if err != nil {
		return err
	}
	outputDriverCfg, err := cl.Get(ctx, root+"/output/config", nil)
	if err != nil {
		return err
	}

	var input config.Input
	var output config.Output
	input.Driver = inputDriverResponse.Node.Value
	output.Driver = outputDriverResponse.Node.Value

	err = json.Unmarshal([]byte(inputDriverCfg.Node.Value), &input.Config)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(outputDriverCfg.Node.Value), &output.Config)
	if err != nil {
		return err
	}

	filtersResponse, err := cl.Get(ctx, root+"/filters", &client.GetOptions{Recursive: true})
	if err != nil {
		return err
	}

	for idx, node := range filtersResponse.Node.Nodes {
		var filter config.Filter
		for _, n := range node.Nodes {
			if n.Key == root+"/filters/"+strconv.Itoa(idx)+"/driver" {
				filter.Driver = n.Value
			} else if n.Key == root+"/filters/"+strconv.Itoa(idx)+"/config" {
				err = json.Unmarshal([]byte(n.Value), &filter.Config)
				if err != nil {
					return err
				}
			} else {
				log.Log(ctx).Warn("Unexpected node", "node", n)
			}
		}
		cfg.Filters = append(cfg.Filters, filter)
	}

	cfg.Input = input
	cfg.Output = output

	return nil

}