package rss

import (
	"bytes"
	"encoding/xml"
	"html"
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
	// a different URL (migration hint). Podcast apps may refuse to subscribe to the proxy
	// feed if the tag indicates the feed lives elsewhere.
	RewriteNewFeedURL string
}

func xmlEscapeText(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

// normalizeTitleBytes makes the byte-extracted <title> comparable with xml.Unmarshal results.
// - strips optional <![CDATA[...]]> wrapper (regex extraction keeps it)
// - unescapes XML/HTML entities (&amp;, &apos;, &#...)
// - trims spaces
func normalizeTitleBytes(b []byte) string {
	s := strings.TrimSpace(string(b))

	// Handle CDATA wrapper extracted by regex: <![CDATA[...]]>
	if strings.HasPrefix(s, "<![CDATA[") && strings.HasSuffix(s, "]]>") {
		s = strings.TrimPrefix(s, "<![CDATA[")
		s = strings.TrimSuffix(s, "]]>")
		s = strings.TrimSpace(s)
	}

	// Decode XML entities to match xml.Unmarshal output
	s = html.UnescapeString(s)

	return strings.TrimSpace(s)
}

// FilterXML filters the original RSS XML by keeping only items whose <title> matches keepTitles.
//
// For advanced behaviors, use FilterXMLWithOptions.
func FilterXML(raw []byte, keepTitles map[string]bool) ([]byte, error) {
	return FilterXMLWithOptions(raw, keepTitles, FilterXMLOptions{})
}

// FilterXMLWithOptions is the same as FilterXML, but allows additional RSS-safe rewrites.
func FilterXMLWithOptions(raw []byte, keepTitles map[string]bool, opts FilterXMLOptions) ([]byte, error) {
	// Optional: rewrite itunes:new-feed-url (kept as byte-level replacement; no XML re-serialization).
	if u := strings.TrimSpace(opts.RewriteNewFeedURL); u != "" {
		if reNewFeedURL.Match(raw) {
			repl := []byte("<itunes:new-feed-url>" + xmlEscapeText(u) + "</itunes:new-feed-url>")
			raw = reNewFeedURL.ReplaceAll(raw, repl)
		}
	}

	matches := reItem.FindAllIndex(raw, -1)
	if len(matches) == 0 {
		// No <item> found; return original.
		return raw, nil
	}

	var buf bytes.Buffer
	last := 0

	for _, m := range matches {
		start, end := m[0], m[1]

		// Write everything between previous item and this item (channel metadata etc.)
		buf.Write(raw[last:start])

		item := raw[start:end]
		titleMatch := reTitle.FindSubmatch(item)
		if titleMatch != nil {
			title := normalizeTitleBytes(titleMatch[1])
			if keepTitles[title] {
				buf.Write(item)
			}
		}

		last = end
	}

	// Write the tail (closing channel/rss tags)
	buf.Write(raw[last:])
	return buf.Bytes(), nil
}
