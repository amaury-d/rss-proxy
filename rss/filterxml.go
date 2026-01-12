package rss

import (
	"bytes"
	"encoding/xml"
	"regexp"
	"strings"
)

// Item-level byte filtering.
// This avoids XML re-serialization issues (namespaces, CDATA, Apple Podcasts compatibility).
var (
	reItem       = regexp.MustCompile(`(?s)<item\b.*?</item>`)
	reTitle      = regexp.MustCompile(`(?s)<title>(.*?)</title>`)
	reNewFeedURL = regexp.MustCompile(`(?is)<itunes:new-feed-url\b[^>]*>.*?</itunes:new-feed-url>`)
)

// FilterXMLOptions controls optional post-processing done on the raw feed.
//
// Note: filtering remains byte-level (we do not re-serialize XML) to preserve namespaces,
// CDATA blocks and Apple Podcasts compatibility.
type FilterXMLOptions struct {
	// RewriteNewFeedURL, when non-empty, rewrites any <itunes:new-feed-url> tag value
	// to the provided URL.
	//
	// This is useful when the upstream feed contains an itunes:new-feed-url pointing to
	// a different URL (migration hint). Some podcast apps (incl. Apple Podcasts) may
	// refuse to subscribe to the proxy feed if the tag indicates the feed lives elsewhere.
	RewriteNewFeedURL string
}

func xmlEscapeText(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

func FilterXML(raw []byte, keepTitles map[string]bool) ([]byte, error) {
	return FilterXMLWithOptions(raw, keepTitles, FilterXMLOptions{})
}

func FilterXMLWithOptions(raw []byte, keepTitles map[string]bool, opts FilterXMLOptions) ([]byte, error) {
	if u := strings.TrimSpace(opts.RewriteNewFeedURL); u != "" {
		if reNewFeedURL.Match(raw) {
			repl := []byte("<itunes:new-feed-url>" + xmlEscapeText(u) + "</itunes:new-feed-url>")
			raw = reNewFeedURL.ReplaceAll(raw, repl)
		}
	}

	matches := reItem.FindAllIndex(raw, -1)
	if len(matches) == 0 {
		return raw, nil
	}

	var buf bytes.Buffer
	last := 0

	for _, m := range matches {
		start, end := m[0], m[1]

		buf.Write(raw[last:start])

		item := raw[start:end]
		titleMatch := reTitle.FindSubmatch(item)

		if titleMatch != nil {
			if keepTitles[string(bytes.TrimSpace(titleMatch[1]))] {
				buf.Write(item)
			}
		}

		last = end
	}

	buf.Write(raw[last:])
	return buf.Bytes(), nil
}
