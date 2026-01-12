package config

import (
	"os"
	"testing"
)

func TestLoadParsesServerBaseURL(t *testing.T) {
	yaml := `
server:
  base_url: https://podcasts.decre.me/rss
feeds:
  - id: bible-en-un-an
    source: https://feeds.example.com/feed.rss
    rules:
      - type: episode_number_min
        min: 640
`

	f, err := os.CreateTemp("", "rss-proxy-config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(yaml); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	_ = f.Close()

	cfg := Load(f.Name())
	if cfg.Server.BaseURL != "https://podcasts.decre.me/rss" {
		t.Fatalf("expected base_url to be parsed, got %q", cfg.Server.BaseURL)
	}
	if len(cfg.Feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(cfg.Feeds))
	}
	if cfg.Feeds[0].ID != "bible-en-un-an" {
		t.Fatalf("unexpected feed id: %q", cfg.Feeds[0].ID)
	}
}
