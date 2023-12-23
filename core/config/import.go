package config

import (
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Servers       []string `yaml:"servers"`
	Hosts         []string `yaml:"hosts"`
	Configuration struct {
		Path     string `yaml:"path"`
		TimeFile string `yaml:"timefile"`
		UPFile   string `yaml:"upfile"`
		Interval int    `yaml:"interval"`
		Write    bool   `yaml:"writetofiles"`
	} `yaml:"configuration"`
}

// FromYAML reads YAML from reader and unmarshals it to Config.
func FromYAML(r io.Reader) (*Config, error) {
	c := &Config{}
	err := yaml.NewDecoder(r).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	return c, nil
}
