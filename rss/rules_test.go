package rss

import (
	"testing"

	"rss-proxy/config"
)

func TestEpisodeNumberMinRule(t *testing.T) {
	feed := RSS{
		Channel: Channel{
			Items: []Item{
				{Title: "Jour 639", Episode: 639},
				{Title: "Jour 640", Episode: 640},
			},
		},
	}

	rules := []config.Rule{
		{Type: "episode_number_min", Min: 640},
	}

	out := ApplyRules(feed, rules)

	if len(out.Channel.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out.Channel.Items))
	}

	if out.Channel.Items[0].Episode != 640 {
		t.Fatal("wrong episode kept")
	}
}

func TestTitleContainsRule(t *testing.T) {
	feed := RSS{
		Channel: Channel{
			Items: []Item{
				{Title: "Normal episode"},
				{Title: "[REDIFF] Old episode"},
			},
		},
	}

	rules := []config.Rule{
		{Type: "title_contains", Value: "[REDIFF]"},
	}

	out := ApplyRules(feed, rules)

	if len(out.Channel.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out.Channel.Items))
	}

	if out.Channel.Items[0].Title != "[REDIFF] Old episode" {
		t.Fatal("wrong item kept")
	}
}

func TestTitleExcludeRule(t *testing.T) {
	feed := RSS{
		Channel: Channel{
			Items: []Item{
				{Title: "Normal episode"},
				{Title: "[REDIFF] Old episode"},
			},
		},
	}

	rules := []config.Rule{
		{Type: "title_excludes", Value: "REDIFF"},
	}

	out := ApplyRules(feed, rules)

	if len(out.Channel.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out.Channel.Items))
	}

	if out.Channel.Items[0].Title != "Normal episode" {
		t.Fatal("wrong item kept")
	}
}

func TestTitleFractionEqualsRule(t *testing.T) {
	feed := RSS{
		Channel: Channel{
			Items: []Item{
				{Title: "[1/2] L'affaire des petits pains au chocolat"},
				{Title: "[2/2] L'affaire des petits pains au chocolat"},
				{Title: "[3/3] Autre affaire"},
				{Title: "[2/3] Encore autre chose"},
				{Title: "Épisode sans fraction"},
			},
		},
	}

	rules := []config.Rule{
		{Type: "title_fraction_equals"},
	}

	out := ApplyRules(feed, rules)

	if len(out.Channel.Items) != 2 {
		t.Fatalf("expected 2 items kept, got %d", len(out.Channel.Items))
	}

	if out.Channel.Items[0].Title != "[2/2] L'affaire des petits pains au chocolat" {
		t.Fatal("wrong first item kept")
	}

	if out.Channel.Items[1].Title != "[3/3] Autre affaire" {
		t.Fatal("wrong second item kept")
	}
}
