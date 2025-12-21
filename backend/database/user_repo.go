package database

import (
	"backend/models"
	"backend/utils"
	"database/sql"
	"time"
)

func CreateUser(db *sql.DB, user *models.User) error {

	hash, hashingErr := utils.HashingPassword(user.Password)

	if hashingErr != nil {
		return hashingErr
	}

	user.PasswordHash = hash
	user.CreatedAt = time.Now()
	user.LastActive = time.Now()

	query := `
	INSERT INTO users (
		username,
		password_hash,
		created_at,
		last_active
	)
	VALUES (?, ?, ?, ?);
	`
	_, err := db.Exec(
		query,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
		user.LastActive,
	)

	user.Password = ""
	user.PasswordHash = ""

	return err
}
