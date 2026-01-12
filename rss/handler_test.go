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
			// Should keep KEEP ME and CDATA KEEP, but not Fish & Chips nor DROP ME.
			{Type: "title_contains", Value: "KEEP"},
		},
	}

	handler := NewHandlerWithBaseURL(feed, cache, "https://podcasts.decre.me/rss")

	req := httptest.NewRequest("GET", "/rss/test.xml", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", w.Code)
	}

	body := w.Body.String()

	if !strings.Contains(body, "<title>KEEP ME</title>") {
		t.Fatal("expected KEEP ME item missing")
	}
	// Critical: ensure CDATA titles are not dropped by the byte-level filter.
	if !strings.Contains(body, "<title><![CDATA[CDATA KEEP]]></title>") {
		t.Fatal("expected CDATA KEEP item missing")
	}

	if strings.Contains(body, "Fish &amp; Chips") {
		t.Fatal("unexpected Fish & Chips item present")
	}

	if strings.Contains(body, "<title>DROP ME</title>") {
		t.Fatal("unexpected DROP ME item present")
	}

	// Ensure itunes:new-feed-url is rewritten to the proxy URL (not upstream).
	if !strings.Contains(body, "<itunes:new-feed-url>https://podcasts.decre.me/rss/test.xml</itunes:new-feed-url>") {
		t.Fatal("expected itunes:new-feed-url to be rewritten to proxy URL")
	}
}
