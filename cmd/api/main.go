package main

import (
	"fmt"
	"log"
	"time"

	"github.com/llan0/go-social/internal/db"
	"github.com/llan0/go-social/internal/env" // can use pkgs like godotenv (this is my own implimentation)
	"github.com/llan0/go-social/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdelConns: env.GetInt("DB_MAX_IDEL_CONN", 30),
			maxIdelTime:  env.GetDuration("DB_MAX_IDEL_TIME", 15*time.Minute),
		},
		env: env.GetString("ENV", "development"),
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdelConns, cfg.db.maxIdelTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	fmt.Println("db connection pool established!")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
