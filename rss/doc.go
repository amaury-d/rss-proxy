// Package rss implements the core logic of the RSS proxy.
//
// It is responsible for:
//   - Fetching RSS feeds with HTTP caching (ETag / If-Modified-Since)
//   - Parsing feeds for rule evaluation (read-only)
//   - Filtering RSS items without XML re-serialization
//   - Preserving full compatibility with podcast clients such as Apple Podcasts
//
// The package intentionally avoids transforming or rewriting audio enclosures
// and focuses on deterministic, item-level filtering.
package rss
