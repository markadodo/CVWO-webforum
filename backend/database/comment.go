package database

import (
	"backend/models"
	"database/sql"
	"errors"
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

func GetCommentOwnerByID(db *sql.DB, commentID int64) (int64, error) {
	commentData, err := ReadCommentByID(db, commentID)

	if err != nil {
		return 0, err
	}

	if commentData == nil {
		return 0, nil
	}

	return commentData.CreatedBy, err
}

func ReadCommentByPostID(db *sql.DB, postID int64, limit int, offset int, sortBy string, order string) ([]models.Comment, error) {
	var comments []models.Comment

	query := `
	SELECT * FROM comments
	WHERE post_id = ?
	ORDER BY ` + sortBy + " " + order + `
	LIMIT ? OFFSET ?`

	rows, err := db.Query(
		query,
		postID,
		limit,
		offset,
	)

	if err != nil {
		return comments, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment models.Comment

		if err := rows.Scan(&comment.ID, &comment.Description, &comment.Likes, &comment.Dislikes, &comment.IsEdited, &comment.PostID, &comment.ParentCommentID, &comment.CreatedBy, &comment.CreatedAt); err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return comments, err
	}

	return comments, nil
}

var ErrDuplicateCommentReaction = errors.New("reaction already exists")

func CreateCommentReaction(db *sql.DB, input *models.CommentReaction) error {
	query := `
	INSERT INTO comments_reactions (
		comment_id,
		user_id,
		reaction
	)
	VALUES (?, ?, ?);
	`
	_, err := db.Exec(query, input.CommentID, input.UserID, input.Reaction)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicateCommentReaction
		}
		return err
	}

	return nil
}

func DeleteCommentReactionByCommentIDAndUserID(db *sql.DB, commentID int64, userID int64) (bool, error) {
	query := "DELETE FROM comments_reactions WHERE comment_id = ? AND user_id = ?"
	res, err := db.Exec(query, commentID, userID)

	if err != nil {
		return false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return true, nil
	}

	return false, nil
}
