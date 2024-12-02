package main

import (
	"github.com/szymczykkrzysztof/social/internal/db"
	"github.com/szymczykkrzysztof/social/internal/env"
	"github.com/szymczykkrzysztof/social/internal/store"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:/socialnetwork?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store)
}
