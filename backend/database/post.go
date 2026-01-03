package database

import (
	"backend/models"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

func CreatePost(db *sql.DB, post *models.Post) error {

	post.CreatedAt = time.Now()

	query := `
	INSERT INTO posts (
		title,
		description,
		topic_id,
		created_by,
		created_at
	)
	VALUES (?, ?, ?, ?, ?);
	`
	result, err := db.Exec(
		query,
		post.Title,
		post.Description,
		post.TopicID,
		post.CreatedBy,
		post.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	post.ID = id

	return nil
}

func ReadPostByID(db *sql.DB, id int64) (*models.Post, error) {
	post := models.Post{}

	query := `
	SELECT id, title, description, topic_id, likes, dislikes, is_edited, views, popularity, created_by, created_at
	FROM posts
	WHERE id = ?
	`
	err := db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Description, &post.TopicID, &post.Likes, &post.Dislikes, &post.IsEdited, &post.Views, &post.Popularity, &post.CreatedBy, &post.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func UpdatePostByID(db *sql.DB, id int64, input *models.UpdatePostInput) (bool, bool, error) {
	updates := []string{}
	args := []interface{}{}

	if input.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *input.Title)
	}

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

	if input.Views != nil {
		updates = append(updates, "views = ?")
		args = append(args, *input.Views)
	}

	if input.Popularity != nil {
		updates = append(updates, "popularity = ?")
		args = append(args, *input.Popularity)
	}

	if len(updates) == 0 {
		return true, false, nil
	}

	query := "UPDATE posts SET " + strings.Join(updates, ", ") + " WHERE id = ?"
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

func DeletePostByID(db *sql.DB, id int64) (bool, error) {
	query := "DELETE FROM posts WHERE id = ?"
	res, err := db.Exec(query, id)

	if err != nil {
		return false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return true, nil
	}

	return false, nil
}

func GetPostOwnerByID(db *sql.DB, postID int64) (int64, error) {
	postData, err := ReadPostByID(db, postID)

	if err != nil {
		return 0, err
	}

	if postData == nil {
		return 0, nil
	}

	return postData.CreatedBy, err
}

func ReadPostByTopicID(db *sql.DB, topicID int64, limit int, offset int, sortBy string, order string) ([]models.Post, error) {
	var posts []models.Post

	query := `
	SELECT * FROM posts
	WHERE topic_id = ?
	ORDER BY ` + sortBy + " " + order + `
	LIMIT ? OFFSET ?`

	rows, err := db.Query(
		query,
		topicID,
		limit,
		offset,
	)

	if err != nil {
		return posts, err
	}

	defer rows.Close()

	for rows.Next() {
		var post models.Post

		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.TopicID, &post.Likes, &post.Dislikes, &post.IsEdited, &post.Views, &post.Popularity, &post.CreatedBy, &post.CreatedAt); err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil
}

func ReadPostBySearchQuery(db *sql.DB, topicID int64, limit int, offset int, sortBy string, order string, searchQuery string) ([]models.Post, error) {
	var posts []models.Post
	args := []interface{}{searchQuery}

	query := `
	SELECT posts.*
	FROM posts
	JOIN posts_fts ON posts.id = posts_fts.rowid
	WHERE posts_fts MATCH ?
	`

	if topicID != 0 {
		query = query + " AND posts.topic_id = ?"
		args = append(args, topicID)
	}

	if sortBy == "relevance" {
		sortBy = "bm25(posts_fts)"
	}

	query = query + " ORDER BY " + sortBy + " " + order + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)

	if err != nil {
		return posts, err
	}

	defer rows.Close()

	for rows.Next() {
		var post models.Post

		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.TopicID, &post.Likes, &post.Dislikes, &post.IsEdited, &post.Views, &post.Popularity, &post.CreatedBy, &post.CreatedAt); err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil
}

var ErrDuplicatePostReaction = errors.New("reaction already exists")

func CreatePostReaction(db *sql.DB, input *models.PostReaction) error {
	query := `
	INSERT INTO posts_reactions (
		post_id,
		user_id,
		reaction
	)
	VALUES (?, ?, ?);
	`

	_, err := db.Exec(query, input.PostID, input.UserID, input.Reaction)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrDuplicatePostReaction
		}
		log.Println(err)
		return err
	}

	return nil
}

func DeletePostReactionByPostIDAndUserID(db *sql.DB, postID int64, userID int64) (bool, error) {
	query := "DELETE FROM posts_reactions WHERE post_id = ? AND user_id = ?"
	res, err := db.Exec(query, postID, userID)

	if err != nil {
		return false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return true, nil
	}

	return false, nil
}
