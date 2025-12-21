package database

import (
	"database/sql"
)

// database creation
func InitDB(db *sql.DB) error {
	usertable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
 		password_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		last_active DATETIME NOT NULL
		);
	`
	topictable := `
	CREATE TABLE IF NOT EXISTS topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		created_by INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		tags TEXT,
		FOREIGN KEY (created_by) REFERENCES users(id)
		);
	`
	posttable := `
	CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		topic_id INTEGER NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		is_edited INTEGER NOT NULL DEFAULT 0,
		views INTEGER NOT NULL DEFAULT 0,
		popularity INTEGER NOT NULL DEFAULT 0,
		created_by INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		tags TEXT,
		FOREIGN KEY (topic_id) REFERENCES topics(id),
		FOREIGN KEY (created_by) REFERENCES users(id)
        );
    `
	//attachment

	commenttable := `
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		is_edited INTEGER NOT NULL DEFAULT 0,
		post_id INTEGER NOT NULL,
		parent_comment_id INTEGER,
		created_by INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (created_by) REFERENCES users(id),
		FOREIGN KEY (parent_comment_id) REFERENCES comments(id)
		);
	`

	tables := []string{
		usertable,
		topictable,
		posttable,
		commenttable}

	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
