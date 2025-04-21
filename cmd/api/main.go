package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lunatictiol/go-based-social-media/internal/env"
	"github.com/lunatictiol/go-based-social-media/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("PORT", ":8080"),
	}
	store := store.NewStorage(nil)
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
