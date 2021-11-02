package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	//timeout := time.Second * 5

	mux.HandleFunc("/some/path/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(timeout)
		w.WriteHeader(http.StatusOK)
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
