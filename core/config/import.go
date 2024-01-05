package config

import (
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Servers       []TargetConfig `yaml:"servers"`

	Hosts         []string `yaml:"hosts"`
	Configuration struct {
		Path       string `yaml:"path"`
		MetricFile string `yaml:"filename"`
		Interval   int    `yaml:"interval"`
		Write      bool   `yaml:"externalcollector"`
		DNSTimeout int    `yaml:"timeout"`
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

func (cfg *Config) TargetConfigByAddr(addr string) TargetConfig {
	for _, t := range cfg.Servers {
		if t.Addr == addr {
			return t
		}
	}

	return TargetConfig{Addr: addr}
}