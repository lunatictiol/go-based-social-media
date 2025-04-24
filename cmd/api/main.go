package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/lunatictiol/go-based-social-media/internal/db"
	"github.com/lunatictiol/go-based-social-media/internal/env"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

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
		addr: env.GetString("PORT", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres"),
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
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
