package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
}

func (a *application) mount() http.Handler {
	r := chi.NewRouter()

	//middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", a.healthCheckHandler)

	})

	return r
}

func (a *application) run(mux http.Handler) error {
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
