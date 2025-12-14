package rss

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// CacheStatus represents the outcome of a cache fetch.
type CacheStatus string

const (
	CacheHit         CacheStatus = "HIT"
	CacheMiss        CacheStatus = "MISS"
	CacheRevalidated CacheStatus = "REVALIDATED"
	CacheStale       CacheStatus = "STALE"
)

type cachedEntry struct {
	body         []byte
	etag         string
	lastModified string
	fetchedAt    time.Time
}

// HTTPCache is a simple in-memory HTTP cache with ETag and If-Modified-Since support.
type HTTPCache struct {
	ttl time.Duration

	mu    sync.RWMutex
	items map[string]*cachedEntry

	now    func() time.Time
	client *http.Client
}

// NewHTTPCache creates a new HTTP cache with a given TTL.
func NewHTTPCache(ttl time.Duration) *HTTPCache {
	return &HTTPCache{
		ttl:    ttl,
		items:  make(map[string]*cachedEntry),
		now:    time.Now,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Fetch retrieves a URL using conditional HTTP requests when possible.
func (c *HTTPCache) Fetch(url string) ([]byte, CacheStatus, error) {
	c.mu.RLock()
	e := c.items[url]
	if e != nil && c.now().Sub(e.fetchedAt) < c.ttl {
		body := e.body
		c.mu.RUnlock()
		return body, CacheHit, nil
	}
	c.mu.RUnlock()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "rss-proxy/1.0")

	if e != nil {
		if e.etag != "" {
			req.Header.Set("If-None-Match", e.etag)
		}
		if e.lastModified != "" {
			req.Header.Set("If-Modified-Since", e.lastModified)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if e != nil && len(e.body) > 0 {
			return e.body, CacheStale, nil
		}
		return nil, "", err
	}

	// Properly handle Close() error for errcheck compliance
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			Logger.Warn("failed to close response body",
				"url", url,
				"error", cerr,
			)
		}
	}()

	switch resp.StatusCode {
	case http.StatusNotModified:
		if e == nil || len(e.body) == 0 {
			return nil, "", fmt.Errorf("received 304 without cached body")
		}

		c.mu.Lock()
		e.fetchedAt = c.now()
		c.mu.Unlock()

		return e.body, CacheRevalidated, nil

	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if e != nil && len(e.body) > 0 {
				return e.body, CacheStale, nil
			}
			return nil, "", err
		}

		c.mu.Lock()
		c.items[url] = &cachedEntry{
			body:         body,
			etag:         resp.Header.Get("ETag"),
			lastModified: resp.Header.Get("Last-Modified"),
			fetchedAt:    c.now(),
		}
		c.mu.Unlock()

		return body, CacheMiss, nil

	default:
		if e != nil && len(e.body) > 0 {
			return e.body, CacheStale, nil
		}
		return nil, "", fmt.Errorf("upstream returned status %d", resp.StatusCode)
	}
}
