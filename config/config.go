// Package config defines the YAML configuration model for rss-proxy.
//
// It provides:
//   - The feed configuration schema
//   - Declarative filtering rules
//   - Minimal validation helpers
//
// The configuration is intentionally simple and explicit, and is designed
// to be loaded once at startup.
package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server `yaml:"server"`
	Feeds  []Feed `yaml:"feeds"`
}

type Server struct {
	// BaseURL is the externally reachable base URL of the rss-proxy feed endpoint.
	// Example: https://podcasts.decre.me/rss
	//
	// When set, rss-proxy rewrites <itunes:new-feed-url> (if present in the upstream
	// feed) to point back to the proxied feed URL, instead of instructing podcast apps
	// to migrate to the upstream provider URL.
	BaseURL string `yaml:"base_url"`
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
