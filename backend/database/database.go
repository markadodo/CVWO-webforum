package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
