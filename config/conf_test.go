package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	conf, err := LoadFile("example.json")
	if err != nil {
		t.Errorf("Expected error to be nil, is %s", err)
	}

	if drv := conf.Input.Driver; drv != "tcp" {
		t.Errorf("Expected input driver to be tcp, is %s", drv)
	}

	if drv := conf.Filters[0].Driver; drv != "uuid" {
		t.Errorf("Expected input driver to be uuid, is %s", drv)
	}

	if port := conf.Input.Config["port"].(float64); port != 9999 {
		t.Errorf("Expected input config value port to be 9999, is %f", port)
	}
}
