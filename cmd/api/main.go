package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/lunatictiol/go-based-social-media/internal/auth"
	"github.com/lunatictiol/go-based-social-media/internal/db"
	"github.com/lunatictiol/go-based-social-media/internal/env"
	"github.com/lunatictiol/go-based-social-media/internal/mailer"
	"github.com/lunatictiol/go-based-social-media/internal/ratelimiter"
	"github.com/lunatictiol/go-based-social-media/internal/store"
	"github.com/lunatictiol/go-based-social-media/internal/store/cache"
	"go.uber.org/zap"
)

var version = ""

//	@title			Go based social media
//	@description	API for Go social medias, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					api/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	err := godotenv.Load()
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	maxOpenConns, err := env.GetInt("DB_MAX_OPEN_CON", 20)

	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	maxIdleConns, err := env.GetInt("DB_MAX_IDLE_CON", 20)

	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	redisDB, err := env.GetInt("REDIS_DB", 0)
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	RATELIMITER_REQUESTS_COUNT, err := env.GetInt("RATELIMITER_REQUESTS_COUNT", 20)
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:        env.GetString("PORT", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres"),
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			apiKey:    env.GetString("MAIL_APIKEY", "apikey"),
			fromEmail: env.GetString("FROM_EMAIL", "from-email"),
		},
		auth: authConfig{
			basic: basicConfig{
				admin:         env.GetString("ADMIN_USER", "admin"),
				adminPassword: env.GetString("ADMIN_PASSWORD", "password"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "gosocialmedia",
			},
		},
		redisConfig: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      redisDB,
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: RATELIMITER_REQUESTS_COUNT,
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}
	logger.Info("connecting to database")
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()
	store := store.NewStorage(db)
	//sendgrid
	//mailer := mailer.NewMailer(cfg.mail.apiKey, cfg.mail.fromEmail)
	mailer, err := mailer.NewMailTrapClient(cfg.mail.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatalf("Error creating mailer file :%v", err)
	}

	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)
	//cache
	var rdb *redis.Client
	if cfg.redisConfig.enabled {
		rdb = cache.NewRedisClient(cfg.redisConfig.addr, cfg.redisConfig.pw, cfg.redisConfig.db)
		logger.Info("redis cache connection established")

		defer rdb.Close()
	}
	cacheStorage := cache.NewRedisStorage(rdb)

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)
	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		ratelimiter:   rateLimiter,
	}
	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
