package database

import (
	"backend/models"
	"backend/utils"
	"database/sql"
	"strings"
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
	result, err := db.Exec(
		query,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
		user.LastActive,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	user.ID = id
	user.Password = ""
	user.PasswordHash = ""

	return nil
}

func ReadUserByID(db *sql.DB, id int64) (*models.User, error) {
	user := models.User{}

	query := `
	SELECT id, username, password_hash, created_at, last_active
	FROM users
	WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.LastActive)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func ReadUserByUsername(db *sql.DB, username string) (*models.User, error) {
	user := models.User{}

	query := `
	SELECT id, username, password_hash, created_at, last_active
	FROM users
	WHERE username = ?
	`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.LastActive)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserByID(db *sql.DB, id int64, input *models.UpdateUserInput) (bool, bool, error) {
	updates := []string{}
	args := []interface{}{}

	if input.Username != nil {
		updates = append(updates, "username = ?")
		args = append(args, *input.Username)
	}

	if input.LastActive != nil {
		updates = append(updates, "last_active = ?")
		args = append(args, *input.LastActive)
	}

	if input.Password != nil {

		hash, err := utils.HashingPassword(*input.Password)
		if err != nil {
			return false, false, err
		}
		updates = append(updates, "password_hash = ?")
		args = append(args, hash)
	}

	if len(updates) == 0 {
		return true, false, nil
	}

	query := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)
	res, err := db.Exec(query, args...)

	if err != nil {
		return false, false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return false, true, nil
	}

	return false, false, nil
}

func DeleteUserByID(db *sql.DB, id int64) (bool, error) {
	query := "DELETE FROM users WHERE id = ?"
	res, err := db.Exec(query, id)

	if err != nil {
		return false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return true, nil
	}

	return false, nil
}

func GetUserOwnerByID(db *sql.DB, userID int64) (int64, error) {
	return userID, nil
}
