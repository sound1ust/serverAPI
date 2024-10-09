package repo

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	log.Println("Connecting to DB...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
		return nil, err
	}
	return db, nil
}
