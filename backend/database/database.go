package database

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// establish connection to the database
func ConnectDB() (*sql.DB, error) {
	connection := os.Getenv("DATABASE_URL_CVWO")
	db, err := sql.Open("pgx", connection)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
