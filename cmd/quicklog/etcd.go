package main

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"github.com/sheenobu/golibs/apps"
	"github.com/sheenobu/golibs/log"

	"strings"
)

func startEtcdQuicklog(mainCtx context.Context, app *apps.App) {

	log.Log(mainCtx).Debug("Connecting to endpoint", "endpoints", strings.Split(etcdEndpoints, ","))
	cfg := client.Config{
		Endpoints: strings.Split(etcdEndpoints, ","),
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Log(mainCtx).Crit("Error connecting to etcd server", "error", err)
		return
	}
	root := "/quicklog/" + instanceName

	log.Log(mainCtx).Debug("Listening for etcd keys", "key", root)

	kapi := client.NewKeysAPI(c)

	app.SpawnSimple("etcd", func(ctx context.Context) {

		w := kapi.Watcher(root, &client.WatcherOptions{
			Recursive: true,
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

				//TODO: (re)load config
			}
		}
	})
}
