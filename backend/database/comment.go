package database

import (
	"backend/models"
	"database/sql"
	"strings"
	"time"
)

func CreateComment(db *sql.DB, comment *models.Comment) error {

	comment.CreatedAt = time.Now()

	query := `
	INSERT INTO comments (
		description,
		post_id,
		parent_comment_id,
		created_by,
		created_at
	)
	VALUES (?, ?, ?, ?, ?);
	`
	result, err := db.Exec(
		query,
		comment.Description,
		comment.PostID,
		comment.ParentCommentID,
		comment.CreatedBy,
		comment.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	comment.ID = id

	return nil
}

func ReadCommentByID(db *sql.DB, id int64) (*models.Comment, error) {
	comment := models.Comment{}

	query := `
	SELECT id, description, likes, dislikes, is_edited, post_id, parent_comment_id, created_by, created_at
	FROM comments
	WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(&comment.ID, &comment.Description, &comment.Likes, &comment.Dislikes, &comment.IsEdited, &comment.PostID, &comment.ParentCommentID, &comment.CreatedBy, &comment.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func UpdateCommentByID(db *sql.DB, id int64, input *models.UpdateCommentInput) (bool, bool, error) {
	updates := []string{}
	args := []interface{}{}

	if input.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *input.Description)
	}

	if input.Likes != nil {
		updates = append(updates, "likes = ?")
		args = append(args, *input.Likes)
	}

	if input.Dislikes != nil {
		updates = append(updates, "dislikes = ?")
		args = append(args, *input.Dislikes)
	}

	if input.IsEdited != nil {
		updates = append(updates, "is_edited = ?")
		args = append(args, *input.IsEdited)
	}

	if len(updates) == 0 {
		return true, false, nil
	}

	query := "UPDATE comments SET " + strings.Join(updates, ", ") + " WHERE id = ?"
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

func DeleteCommentByID(db *sql.DB, id int64) (bool, error) {
	query := "DELETE FROM comments WHERE id = ?"
	res, err := db.Exec(query, id)

	if err != nil {
		return false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return true, nil
	}

	return false, nil
}
