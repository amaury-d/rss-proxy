package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Feeds []Feed `yaml:"feeds"`
}

type Feed struct {
	ID     string `yaml:"id"`
	Source string `yaml:"source"`
	Rules  []Rule `yaml:"rules"`
}

type Rule struct {
	Type  string `yaml:"type"`
	Min   int    `yaml:"min,omitempty"`
	Value string `yaml:"value,omitempty"`
}

func Load(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
