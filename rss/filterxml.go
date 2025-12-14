package rss

import (
	"bytes"
	"regexp"
)

// Item-level byte filtering.
// This avoids XML re-serialization issues (namespaces, CDATA, Apple Podcasts compatibility).
var (
	reItem  = regexp.MustCompile(`(?s)<item\b.*?</item>`)
	reTitle = regexp.MustCompile(`(?s)<title>(.*?)</title>`)
)

func FilterXML(raw []byte, keepTitles map[string]bool) ([]byte, error) {
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
