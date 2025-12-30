package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	LastActive   time.Time `json:"last_active"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	Username   *string    `json:"username"`
	Password   *string    `json:"password"`
	LastActive *time.Time `json:"last_active"`
}

type LoginUserData struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
