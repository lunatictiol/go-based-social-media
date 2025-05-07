package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lunatictiol/go-based-social-media/docs"
	"github.com/lunatictiol/go-based-social-media/internal/auth"
	"github.com/lunatictiol/go-based-social-media/internal/mailer"
	"github.com/lunatictiol/go-based-social-media/internal/ratelimiter"
	"github.com/lunatictiol/go-based-social-media/internal/store"
	"github.com/lunatictiol/go-based-social-media/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	ratelimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	apiURL      string
	db          dbConfig
	env         string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisConfig redisConfig
	rateLimiter ratelimiter.Config
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
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

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
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
			r.Post("/register", a.registerUserHandler)
			r.Post("/login", a.loginUserHandler)
		})

		//post handler
		r.Route("/post", func(r chi.Router) {
			r.Use(a.AuthTokenMiddleware)
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
				r.Use(a.AuthTokenMiddleware)

				r.Get("/", a.getUserHandler)
				r.Put("/follow", a.followUserHandler)
				r.Put("/unfollow", a.unfollowUserHandler)
			})

			//feed handler
			r.Group(func(r chi.Router) {
				r.Use(a.AuthTokenMiddleware)
				r.Get("/feed", a.getUserFeedHandler)
			})
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
	srv := &http.Server{
		Addr:         a.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	a.logger.Infow("server has started", "addr", a.config.addr, "env", a.config.env)
	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		a.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	a.logger.Infow("server has started", "addr", a.config.addr, "env", a.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	a.logger.Infow("server has stopped", "addr", a.config.addr, "env", a.config.env)

	return nil
}
