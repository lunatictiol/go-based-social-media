package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (a *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func (a *application) run() error {
	mux := http.NewServeMux()
	// Start the server
	s := &http.Server{
		Addr:         a.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Starting server on %s", s.Addr)
	return s.ListenAndServe()
}
