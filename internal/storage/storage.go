package storage

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func NewStorage(pathDb string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", pathDb)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	sqlComand, err := os.ReadFile("./internal/storage/migrations/001_init.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
		return nil, err
	}

	_, err = db.Exec(string(sqlComand))
	if err != nil {
		log.Fatalf("Failed to execute SQL commands: %v", err)
		return nil, err
	}

	return db, nil
}
