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
