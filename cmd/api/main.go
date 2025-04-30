package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/lunatictiol/go-based-social-media/internal/db"
	"github.com/lunatictiol/go-based-social-media/internal/env"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

const version = "0.0.1"

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
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	maxOpenConns, err := env.GetInt("DB_MAX_OPEN_CON", 20)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	maxIdleConns, err := env.GetInt("DB_MAX_IDLE_CON", 20)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:   env.GetString("PORT", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres"),
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}
	fmt.Println("connecting to database")
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
