package rss

import (
	"regexp"
	"strconv"
	"strings"

	"rss-proxy/config"
)

// Matches patterns like [1/2], [2/2], [10/10]
var reFraction = regexp.MustCompile(`\[(\d+)\s*/\s*(\d+)\]`)

// ApplyRules filters RSS items according to the configured rules.
func ApplyRules(feed RSS, rules []config.Rule) RSS {
	filtered := RSS{
		Channel: Channel{
			Title: feed.Channel.Title,
		},
	}

ITEM:
	for _, item := range feed.Channel.Items {
		for _, rule := range rules {
			if !matchRule(item, rule) {
				continue ITEM
			}
		}
		filtered.Channel.Items = append(filtered.Channel.Items, item)
	}

	return filtered
}

func matchRule(item Item, rule config.Rule) bool {
	title := item.Title

	switch rule.Type {

	case "length_max":
		// Keep items whose iTunes duration is <= the configured max.
		// Supported formats for both item.Duration and rule.Value:
		//   - "SS" (seconds)
		//   - "MM:SS"
		//   - "HH:MM:SS"
		maxSec, ok := parseITunesDurationToSeconds(rule.Value)
		if !ok {
			// Invalid config: ignore the rule (keep the item).
			Logger.Warn("Can't parse duration")
			return true
		}
		if strings.TrimSpace(item.Duration) == "" {
			// Some feeds don't provide duration; in that case, don't drop items.
			Logger.Warn("No duration provided")
			return true
		}
		durSec, ok := parseITunesDurationToSeconds(item.Duration)
		if !ok {
			// Unparseable duration: keep the item.
			Logger.Warn("Can't parse duration")
			return true
		}
		return durSec <= maxSec

	case "title_contains":
		return strings.Contains(strings.ToUpper(title), strings.ToUpper(rule.Value))

	case "title_excludes":
		return !strings.Contains(strings.ToUpper(title), strings.ToUpper(rule.Value))

	case "title_regex":
		re := regexp.MustCompile(rule.Value)
		return re.MatchString(title)

	case "episode_number_min":
		if item.Episode > 0 {
			return item.Episode >= rule.Min
		}

		// Fallback: extract episode number from title
		re := regexp.MustCompile(`\b(\d{3,4})\b`)
		m := re.FindStringSubmatch(title)
		if m == nil {
			return false
		}

		n, _ := strconv.Atoi(m[1])
		return n >= rule.Min

	case "title_fraction_equals":
		// Keep only items where [x/y] and x == y
		m := reFraction.FindStringSubmatch(title)
		if m == nil {
			return false
		}

		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])

		return x == y

	default:
		return true
	}
}

// parseITunesDurationToSeconds parses common iTunes duration formats.
//
// iTunes duration can be either:
//   - integer seconds ("1234")
//   - "MM:SS"
//   - "HH:MM:SS" (or "H:MM:SS")
//
// Returns (seconds, true) on success, (0, false) otherwise.
func parseITunesDurationToSeconds(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}

	if !strings.Contains(s, ":") {
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 {
			return 0, false
		}
		return n, true
	}

	parts := strings.Split(s, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}

	toInt := func(p string) (int, bool) {
		p = strings.TrimSpace(p)
		if p == "" {
			return 0, false
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return 0, false
		}
		return n, true
	}

	if len(parts) == 2 {
		mm, ok := toInt(parts[0])
		if !ok {
			return 0, false
		}
		ss, ok := toInt(parts[1])
		if !ok {
			return 0, false
		}
		return mm*60 + ss, true
	}

	hh, ok := toInt(parts[0])
	if !ok {
		return 0, false
	}
	mm, ok := toInt(parts[1])
	if !ok {
		return 0, false
	}
	ss, ok := toInt(parts[2])
	if !ok {
		return 0, false
	}
	return hh*3600 + mm*60 + ss, true
}
