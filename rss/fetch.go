package rss

import (
	"io"
	"net/http"
	"time"
)

// Fetch est une variable pour permettre le mock en test
var Fetch = fetch

func fetch(url string) ([]byte, error) {
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
