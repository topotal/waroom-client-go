package main

import (
	"fmt"
	"os"

	"github.com/go-faster/yaml"
)

// Config is the optional .genwaroom.yaml. An absent or empty file means
// pure convention.
type Config struct {
	// Include limits generation to these operationIds. Empty = all.
	Include []string `yaml:"include"`
	// Skip excludes operationIds (kept hand-written or unsupported).
	Skip []string `yaml:"skip"`
	// Groups overrides the spec tag per operationId.
	Groups map[string]string `yaml:"groups"`
	Rename Rename            `yaml:"rename"`
}

type Rename struct {
	// Methods overrides the public method name per operationId.
	Methods map[string]string `yaml:"methods"`
	// Aliases forces the public name of a gen type.
	Aliases map[string]string `yaml:"aliases"`
}

func loadConfig(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return cfg, nil
}

// wants reports whether operationId id should be generated.
func (c *Config) wants(id string) bool {
	for _, s := range c.Skip {
		if s == id {
			return false
		}
	}
	if len(c.Include) == 0 {
		return true
	}
	for _, s := range c.Include {
		if s == id {
			return true
		}
	}
	return false
}
