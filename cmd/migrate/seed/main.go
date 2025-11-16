package main

import (
	"log"
	"time"

	"github.com/llan0/go-social/internal/db"
	"github.com/llan0/go-social/internal/env"
	"github.com/llan0/go-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, time.Minute*15)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
