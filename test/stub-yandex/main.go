package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	layout, err := os.ReadFile("search.html")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(layout); err != nil {
			log.Printf("failed to wite: %s", err)
		}
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("healthy")); err != nil {
			log.Printf("failed to write status healthy: %s", err)
		}
	})

	if err := (&http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      mux,
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
	}).ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
