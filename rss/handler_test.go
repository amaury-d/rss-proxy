package rss

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rss-proxy/config"
)

func TestHandlerEndToEnd(t *testing.T) {
	// Fake upstream
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sampleRSS))
	}))
	defer srv.Close()

	cache := NewHTTPCache(0)
	cache.client = srv.Client()

	feed := config.Feed{
		ID:     "test",
		Source: srv.URL,
		Rules: []config.Rule{
			{Type: "title_contains", Value: "KEEP"},
		},
	}

	// Configure the external base URL used for itunes:new-feed-url rewrite.
	// This matches the config.yml server.base_url behavior.
	handler := NewHandlerWithBaseURL(feed, cache, "https://podcasts.decre.me/rss")

	req := httptest.NewRequest("GET", "/rss/test.xml", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", w.Code)
	}

	body := w.Body.String()

	if !strings.Contains(body, "KEEP ME") {
		t.Fatal("expected item missing in response")
	}

	if strings.Contains(body, "DROP ME") {
		t.Fatal("unexpected item present in response")
	}

	// Ensure itunes:new-feed-url is rewritten to the proxy URL (not upstream).
	if !strings.Contains(body, "<itunes:new-feed-url>https://podcasts.decre.me/rss/test.xml</itunes:new-feed-url>") {
		t.Fatal("expected itunes:new-feed-url to be rewritten to proxy URL")
	}
}
