package database

import (
	"backend/models"
	"database/sql"
	"errors"
	"log"
	"strconv"
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
	VALUES ($1, $2, $3, $4, $5);
	`
	_, err := db.Exec(
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

	return nil
}

func ReadPostByID(db *sql.DB, id int64) (*models.Post, error) {
	post := models.Post{}

	query := `
	SELECT id, title, description, topic_id, likes, dislikes, is_edited, views, popularity, created_by, created_at
	FROM posts
	WHERE id = $1
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
	counter := 1

	if input.Title != nil {
		updates = append(updates, "title = $"+strconv.Itoa(counter))
		args = append(args, *input.Title)
		counter += 1
	}

	if input.Description != nil {
		updates = append(updates, "description = $"+strconv.Itoa(counter))
		args = append(args, *input.Description)
		counter += 1
	}

	if input.Likes != nil {
		updates = append(updates, "likes = $"+strconv.Itoa(counter))
		args = append(args, *input.Likes)
		counter += 1
	}

	if input.Dislikes != nil {
		updates = append(updates, "dislikes = $"+strconv.Itoa(counter))
		args = append(args, *input.Dislikes)
		counter += 1
	}

	if input.IsEdited != nil {
		updates = append(updates, "is_edited = $"+strconv.Itoa(counter))
		args = append(args, *input.IsEdited)
		counter += 1
	}

	if input.Views != nil {
		updates = append(updates, "views = $"+strconv.Itoa(counter))
		args = append(args, *input.Views)
		counter += 1
	}

	if input.Popularity != nil {
		updates = append(updates, "popularity = $"+strconv.Itoa(counter))
		args = append(args, *input.Popularity)
		counter += 1
	}

	if len(updates) == 0 {
		return true, false, nil
	}

	query := "UPDATE posts SET " + strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(counter)
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
	query := "DELETE FROM posts WHERE id = $1"
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
	SELECT id, title, description, topic_id, likes, dislikes, is_edited, views, popularity, created_by, created_at 
	FROM posts
	WHERE topic_id = $1
	ORDER BY ` + sortBy + " " + order + `
	LIMIT $2 OFFSET $3`

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
	counter := 2

	query := `
	SELECT id, title, description, topic_id, likes, dislikes, is_edited, views, popularity, created_by, created_at
	FROM posts, plainto_tsquery('english', $1) AS query
	WHERE document @@ query
	`

	if topicID != 0 {
		query = query + " AND posts.topic_id = $" + strconv.Itoa(counter)
		args = append(args, topicID)
		counter += 1
	}

	if sortBy == "relevance" {
		sortBy = "ts_rank(document, query)"
	}

	query = query + " ORDER BY " + sortBy + " " + order + " LIMIT $" + strconv.Itoa(counter) + " OFFSET $" + strconv.Itoa(counter+1)
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
	VALUES ($1, $2, $3);
	`

	_, err := db.Exec(query, input.PostID, input.UserID, input.Reaction)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrDuplicatePostReaction
		}
		log.Println(err)
		return err
	}

	return nil
}

func DeletePostReactionByPostIDAndUserID(db *sql.DB, postID int64, userID int64) (bool, error) {
	query := "DELETE FROM posts_reactions WHERE post_id = $1 AND user_id = $2"
	res, err := db.Exec(query, postID, userID)

	if err != nil {
		return false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return true, nil
	}

	return false, nil
}

func ReadPost(db *sql.DB, limit int, offset int, sortBy string, order string) ([]models.Post, error) {
	var posts []models.Post
	args := []interface{}{}

	query := `
	SELECT id, title, description, topic_id, likes, dislikes, is_edited, views, popularity, created_by, created_at
	FROM posts
	`

	query = query + " ORDER BY " + sortBy + " " + order + " LIMIT $1 OFFSET $2"
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
