package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/fission/fission/environments/fetcher"
)

// Usage: fetcher <shared volume path>
func main() {
	dir := os.Args[1]
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModeDir|0700)
			if err != nil {
				log.Fatalf("Error creating directory: %v", err)
			}
		}
	}
	fetcherInstance, err := fetcher.MakeFetcher(dir)
	if err != nil {
		log.Fatalf("Error making fetcher: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", fetcherInstance.FetchHandler)
	mux.HandleFunc("/upload", fetcherInstance.UploadHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Infof("Fetcher is ready to receive requests")
	http.ListenAndServe(":8000", mux)
}
