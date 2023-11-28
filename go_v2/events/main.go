package main

import (
	"database/sql"
	"log"

	"github.com/amirul-zafrin/event/api"
	db "github.com/amirul-zafrin/event/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/events?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.New(conn)
	server := api.NewServer(store)

	server.Start(":3000")
}
