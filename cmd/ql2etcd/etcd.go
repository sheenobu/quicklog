package main

import (
	"strconv"

	"encoding/json"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	"github.com/sheenobu/golibs/log"
	"github.com/sheenobu/quicklog/config"

	"strings"
)

func syncToEtcd(ctx context.Context, cfg *config.Config) error {
	log.Log(ctx).Debug("Connecting to etcd")
	etcdCfg := client.Config{
		Endpoints: strings.Split(etcdEndpoints, ","),
	}

	c, err := client.New(etcdCfg)
	if err != nil {
		return err
	}

	root := "/quicklog/" + instanceName

	kapi := client.NewKeysAPI(c)

	input := cfg.Input

	inputConfig, err := json.Marshal(input.Config)
	if err != nil {
		log.Log(ctx).Error("Error converting input config to JSON data", "error", err)
		return err
	}

	output := cfg.Output

	outputConfig, err := json.Marshal(output.Config)
	if err != nil {
		log.Log(ctx).Error("Error converting output config to JSON data", "error", err)
		return err
	}

	filters := cfg.Filters

	kapi.Set(ctx, root+"/input/driver", input.Driver, nil)
	kapi.Set(ctx, root+"/input/config", string(inputConfig), nil)
	kapi.Set(ctx, root+"/output/driver", output.Driver, nil)
	kapi.Set(ctx, root+"/output/config", string(outputConfig), nil)

	for idx, filter := range filters {
		filterConfig, err := json.Marshal(filter.Config)
		if err != nil {
			log.Log(ctx).Error("Error converting filter config to JSON data", "error", err)
			return err
		}

		kapi.Set(ctx, root+"/filters/"+strconv.Itoa(idx)+"/driver", filter.Driver, nil)
		kapi.Set(ctx, root+"/filters/"+strconv.Itoa(idx)+"/config", string(filterConfig), nil)
	}

	kapi.Set(ctx, root+"/reload", "1", nil)

	return nil
}
