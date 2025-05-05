package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lunatictiol/go-based-social-media/docs"
	"github.com/lunatictiol/go-based-social-media/internal/mailer"
	"github.com/lunatictiol/go-based-social-media/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
}

type config struct {
	addr        string
	apiURL      string
	db          dbConfig
	env         string
	mail        mailConfig
	frontendURL string
	auth        auth
}

type auth struct {
	basic basicConfig
}

type basicConfig struct {
	admin         string
	adminPassword string
}
type mailConfig struct {
	exp       time.Duration
	apiKey    string
	fromEmail string
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

//routing

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

	r.Route("/api/v1", func(r chi.Router) {
		r.With(a.basicAuthMiddleware()).Get("/health", a.healthCheckHandler)
		docsURL := fmt.Sprintf("%s/swagger/doc.json", a.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		//auth
		r.Route("/authenticate", func(r chi.Router) {
			r.Post("/user", a.registerUserHandler)
		})

		//post handler
		r.Route("/post", func(r chi.Router) {
			r.Post("/", a.createPosthandler)
			r.Post("/comment", a.createCommentHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(a.postContextMiddleware)
				r.Get("/", a.getPostHandler)
				r.Delete("/", a.deletePostHandler)
				r.Patch("/", a.updatePostHandler)

			})
		})

		//user handler
		r.Route("/user", func(r chi.Router) {

			r.Put("/activate/{token}", a.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(a.userContextMiddleware)
				r.Get("/", a.getUserHandler)
				r.Put("/follow", a.followUserHandler)
				r.Put("/unfollow", a.unfollowUserHandler)
			})
		})

		//feed handler
		r.Group(func(r chi.Router) {
			r.Get("/feed", a.getUserFeedHandler)
		})

	})

	return r
}

func (a *application) run(mux http.Handler) error {
	//docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = a.config.apiURL
	docs.SwaggerInfo.BasePath = "/api/v1"
	// Start the server
	s := &http.Server{
		Addr:         a.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	a.logger.Infow("server has started", "addr", a.config.addr, "env", a.config.env)
	return s.ListenAndServe()
}
