package database

import (
	"backend/models"
	"backend/utils"
	"database/sql"
	"strconv"
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
	VALUES ($1, $2, $3, $4);
	`
	_, err := db.Exec(
		query,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
		user.LastActive,
	)

	if err != nil {
		return err
	}

	user.Password = ""
	user.PasswordHash = ""

	return nil
}

func ReadUserByID(db *sql.DB, id int64) (*models.User, error) {
	user := models.User{}

	query := `
	SELECT id, username, password_hash, created_at, last_active
	FROM users
	WHERE id = $1
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
	WHERE username = $1
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
	counter := 1

	if input.Username != nil {
		placeholder := strconv.Itoa(counter)
		updates = append(updates, "username = $"+placeholder)
		args = append(args, *input.Username)
		counter += 1
	}

	if input.LastActive != nil {
		placeholder := strconv.Itoa(counter)
		updates = append(updates, "last_active = $"+placeholder)
		args = append(args, *input.LastActive)
		counter += 1
	}

	if input.Password != nil {
		placeholder := strconv.Itoa(counter)
		hash, err := utils.HashingPassword(*input.Password)
		if err != nil {
			return false, false, err
		}
		updates = append(updates, "password_hash = $"+placeholder)
		args = append(args, hash)
		counter += 1
	}

	if len(updates) == 0 {
		return true, false, nil
	}

	placeholder := strconv.Itoa(counter)
	query := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = $" + placeholder
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
	query := "DELETE FROM users WHERE id = $1"
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
