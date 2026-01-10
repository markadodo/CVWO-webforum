package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// establish connection to the database
func ConnectDB() (*sql.DB, error) {
	connection := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
