package models

import "time"

type Comment struct {
	ID              int64     `json:"id"`
	Description     string    `json:"description"`
	Likes           int       `json:"likes"`
	Dislikes        int       `json:"dislikes"`
	IsEdited        int       `json:"is_edited"`
	PostID          int64     `json:"post_id"`
	ParentCommentID *int64    `json:"parent_comment_id"`
	CreatedBy       int64     `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	Username        string    `json:"username"`
}

type CreateCommentInput struct {
	Description     string `json:"description"`
	PostID          int64  `json:"post_id"`
	ParentCommentID *int64 `json:"parent_comment_id"`
	CreatedBy       int64  `json:"created_by"`
}

type UpdateCommentInput struct {
	Description *string `json:"description"`
	Likes       *int    `json:"likes"`
	Dislikes    *int    `json:"dislikes"`
	IsEdited    *int    `json:"is_edited"`
}

type CommentReaction struct {
	ID        int64 `json:"id"`
	CommentID int64 `json:"comment_id"`
	UserID    int64 `json:"user_id"`
	Reaction  bool  `json:"reaction"`
}

type CreateCommentReactionInput struct {
	Reaction bool `json:"reaction"`
}
