package main

import (
	"log"
	"net/http"

	"rss-proxy/config"
	"rss-proxy/rss"
)

func main() {
	cfg := config.Load("config.yml")

	for _, feed := range cfg.Feeds {
		handler := rss.NewHandlerWithDefaultCache(feed)
		http.Handle("/rss/"+feed.ID+".xml", handler)
	}

	http.HandleFunc("/rss/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Println("RSS proxy listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
