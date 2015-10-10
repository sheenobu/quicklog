package config

import (
	"encoding/json"
	"io"
	"os"
)

// Config defines a single configuration file
type Config struct {
	Input   Input
	Filters []Filter
	Output  Output
}

// Input defines the input driver and configuration flags passed to the driver
type Input struct {
	Driver string
	Parser string
	Config map[string]interface{}
}

// Filter defines the filter driver and the configuration flags passed to the driver
type Filter struct {
	Driver string
	Config map[string]interface{}
}

// Output defines the output driver and the configuration flags passed to the driver
type Output struct {
	Driver string
	Config map[string]interface{}
}

// LoadFile loads the config via a filepath
func LoadFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return Load(f)
}

// Load laods the config via an io.Reader
func Load(r io.Reader) (*Config, error) {
	var conf Config

	dec := json.NewDecoder(r)
	err := dec.Decode(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
