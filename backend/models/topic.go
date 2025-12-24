package models

import "time"

type Topic struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTopicInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedBy   int64  `json:"created_by"`
}

type UpdateTopicInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
