package rss

import (
	"net/http"
	"strings"
	"time"

	"rss-proxy/config"
)

// Handler serves a filtered RSS feed for a single podcast.
type Handler struct {
	feed  config.Feed
	cache *HTTPCache
	// baseURL is the externally reachable base URL (from config.server.base_url).
	// When set, it is used to rewrite <itunes:new-feed-url> so podcast apps keep
	// subscribers on the proxied feed URL.
	baseURL string
}

// NewHandler creates a handler with an injected HTTP cache.
func NewHandler(feed config.Feed, cache *HTTPCache) http.Handler {
	return &Handler{
		feed:  feed,
		cache: cache,
	}
}

// NewHandlerWithBaseURL creates a handler with an injected HTTP cache and an external base URL.
func NewHandlerWithBaseURL(feed config.Feed, cache *HTTPCache, baseURL string) http.Handler {
	return &Handler{
		feed:    feed,
		cache:   cache,
		baseURL: baseURL,
	}
}

// NewHandlerWithDefaultCache creates a handler with a default in-memory HTTP cache.
func NewHandlerWithDefaultCache(feed config.Feed) http.Handler {
	return NewHandler(feed, NewHTTPCache(15*time.Minute))
}

// NewHandlerWithDefaultCacheAndBaseURL creates a handler with a default in-memory HTTP cache
// and an external base URL.
func NewHandlerWithDefaultCacheAndBaseURL(feed config.Feed, baseURL string) http.Handler {
	return NewHandlerWithBaseURL(feed, NewHTTPCache(15*time.Minute), baseURL)
}

func feedURLFromBase(baseURL, feedID string) string {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return ""
	}
	baseURL = strings.TrimRight(baseURL, "/")
	return baseURL + "/" + feedID + ".xml"
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Logger.Info("fetching feed",
		"feed_id", h.feed.ID,
		"source", h.feed.Source,
	)

	// Fetch RSS using HTTP cache (ETag / If-Modified-Since)
	raw, status, err := h.cache.Fetch(h.feed.Source)
	if err != nil {
		Logger.Error("failed to fetch feed",
			"feed_id", h.feed.ID,
			"source", h.feed.Source,
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	Logger.Info("fetch completed",
		"feed_id", h.feed.ID,
		"cache_status", status,
		"bytes", len(raw),
	)

	// Parse RSS for rule evaluation only (read-only)
	parsed, err := Parse(raw)
	if err != nil {
		Logger.Error("failed to parse feed",
			"feed_id", h.feed.ID,
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Info("parsed feed",
		"feed_id", h.feed.ID,
		"items_total", len(parsed.Channel.Items),
	)

	// Apply filtering rules
	filtered := ApplyRules(parsed, h.feed.Rules)

	Logger.Info("rules applied",
		"feed_id", h.feed.ID,
		"items_kept", len(filtered.Channel.Items),
		"items_dropped", len(parsed.Channel.Items)-len(filtered.Channel.Items),
	)

	// Build allow-list of titles to keep
	keepTitles := make(map[string]bool, len(filtered.Channel.Items))
	for _, item := range filtered.Channel.Items {
		keepTitles[item.Title] = true
	}

	// Filter original XML at item level (byte-for-byte)
	//
	// Also rewrite <itunes:new-feed-url> if configured.
	xmlOut, err := FilterXMLWithOptions(raw, keepTitles, FilterXMLOptions{
		RewriteNewFeedURL: feedURLFromBase(h.baseURL, h.feed.ID),
	})
	if err != nil {
		Logger.Error("failed to filter xml",
			"feed_id", h.feed.ID,
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=900")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(xmlOut); err != nil {
		Logger.Error("failed to write response",
			"feed_id", h.feed.ID,
			"error", err,
		)
	}

	Logger.Info("feed served",
		"feed_id", h.feed.ID,
		"bytes", len(xmlOut),
	)
}
