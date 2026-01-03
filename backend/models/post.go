package models

import "time"

type Post struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TopicID     int64     `json:"topic_id"`
	Likes       int       `json:"likes"`
	Dislikes    int       `json:"dislikes"`
	IsEdited    int       `json:"is_edited"`
	Views       int       `json:"views"`
	Popularity  int       `json:"popularity"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreatePostInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TopicID     int64  `json:"topic_id"`
	CreatedBy   int64  `json:"created_by"`
}

type UpdatePostInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Likes       *int    `json:"likes"`
	Dislikes    *int    `json:"dislikes"`
	IsEdited    *int    `json:"is_edited"`
	Views       *int    `json:"views"`
	Popularity  *int    `json:"popularity"`
}

type PostReaction struct {
	ID       int64 `json:"id"`
	PostID   int64 `json:"post_id"`
	UserID   int64 `json:"user_id"`
	Reaction bool  `json:"reaction"`
}

type CreatePostReactionInput struct {
	Reaction bool `json:"reaction"`
}
