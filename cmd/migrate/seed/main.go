package main

import (
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
	addr := env.GetString("DB_ADDR", "postgres")
	dbs, err := db.New(addr, 3, 3, "5m")
	if err != nil {
		log.Panic(err)
	}
	defer dbs.Close()
	store := store.NewStorage(dbs)
	db.Seed(store)
}
