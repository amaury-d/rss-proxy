package rss

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestHTTPCache_HitMissRevalidate(t *testing.T) {
	var reqCount int32

	bodyV1 := []byte("<rss><channel><title>v1</title></channel></rss>")
	etag := `"v1"`
	lastMod := "Mon, 01 Jan 2024 10:00:00 GMT"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)

		// If client sends conditional GET, answer 304
		if r.Header.Get("If-None-Match") == etag || r.Header.Get("If-Modified-Since") == lastMod {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("ETag", etag)
		w.Header().Set("Last-Modified", lastMod)
		w.WriteHeader(http.StatusOK)
		w.Write(bodyV1)
	}))
	defer srv.Close()

	cache := NewHTTPCache(10 * time.Second)
	// deterministic time
	now := time.Date(2025, 12, 13, 12, 0, 0, 0, time.UTC)
	cache.now = func() time.Time { return now }
	cache.client = srv.Client()

	// 1) MISS
	b, st, err := cache.Fetch(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if st != CacheMiss {
		t.Fatalf("expected MISS, got %s", st)
	}
	if string(b) != string(bodyV1) {
		t.Fatal("unexpected body")
	}
	if atomic.LoadInt32(&reqCount) != 1 {
		t.Fatalf("expected 1 request, got %d", reqCount)
	}

	// 2) HIT within TTL (no network)
	now = now.Add(2 * time.Second)
	b, st, err = cache.Fetch(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if st != CacheHit {
		t.Fatalf("expected HIT, got %s", st)
	}
	if atomic.LoadInt32(&reqCount) != 1 {
		t.Fatalf("expected still 1 request, got %d", reqCount)
	}

	// 3) After TTL -> REVALIDATED (304)
	cache.ttl = 1 * time.Second
	now = now.Add(2 * time.Second)
	b, st, err = cache.Fetch(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if st != CacheRevalidated {
		t.Fatalf("expected REVALIDATED, got %s", st)
	}
	if string(b) != string(bodyV1) {
		t.Fatal("unexpected body after 304")
	}
	if atomic.LoadInt32(&reqCount) != 2 {
		t.Fatalf("expected 2 requests, got %d", reqCount)
	}
}

func TestHTTPCache_StaleOnUpstreamFailure(t *testing.T) {
	// Server that works once, then closes (simulates upstream failure)
	body := []byte("<rss><channel><title>ok</title></channel></rss>")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"ok"`)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	cache := NewHTTPCache(0)
	cache.client = srv.Client()

	// Prime cache
	b, st, err := cache.Fetch(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if st != CacheMiss {
		t.Fatalf("expected MISS, got %s", st)
	}
	if string(b) != string(body) {
		t.Fatal("unexpected body")
	}

	// Kill upstream
	srv.Close()

	// Fetch again => should serve stale
	b, st, err = cache.Fetch(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if st != CacheStale {
		t.Fatalf("expected STALE, got %s", st)
	}
	if string(b) != string(body) {
		t.Fatal("unexpected stale body")
	}
}
