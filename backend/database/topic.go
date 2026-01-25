package database

import (
	"backend/models"
	"database/sql"
	"strconv"
	"strings"
	"time"
)

func CreateTopic(db *sql.DB, topic *models.Topic) error {

	topic.CreatedAt = time.Now()

	query := `
	INSERT INTO topics (
		title,
		description,
		created_by,
		created_at
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`
	err := db.QueryRow(
		query,
		topic.Title,
		topic.Description,
		topic.CreatedBy,
		topic.CreatedAt,
	).Scan(&topic.ID)

	if err != nil {
		return err
	}

	return nil
}

func ReadTopicByID(db *sql.DB, id int64) (*models.Topic, error) {
	topic := models.Topic{}

	query := `
	SELECT id, title, description, created_by, created_at
	FROM topics
	WHERE id = $1
	`
	err := db.QueryRow(query, id).Scan(&topic.ID, &topic.Title, &topic.Description, &topic.CreatedBy, &topic.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &topic, nil
}

func UpdateTopicByID(db *sql.DB, id int64, input *models.UpdateTopicInput) (bool, bool, error) {
	updates := []string{}
	args := []interface{}{}
	counter := 1

	if input.Title != nil {
		placeholder := strconv.Itoa(counter)
		updates = append(updates, "title = $"+placeholder)
		args = append(args, *input.Title)
		counter += 1
	}

	if input.Description != nil {
		placeholder := strconv.Itoa(counter)
		updates = append(updates, "description = $"+placeholder)
		args = append(args, *input.Description)
		counter += 1
	}

	if len(updates) == 0 {
		return true, false, nil
	}
	placeholder := strconv.Itoa(counter)
	query := "UPDATE topics SET " + strings.Join(updates, ", ") + " WHERE id = $" + placeholder
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

func DeleteTopicByID(db *sql.DB, id int64) (bool, error) {
	query := "DELETE FROM topics WHERE id = $1"
	res, err := db.Exec(query, id)

	if err != nil {
		return false, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return true, nil
	}

	return false, nil
}

func GetTopicOwnerByID(db *sql.DB, topicID int64) (int64, error) {
	topicData, err := ReadTopicByID(db, topicID)

	if err != nil {
		return 0, err
	}

	if topicData == nil {
		return 0, nil
	}

	return topicData.CreatedBy, err
}

func ReadTopicBySearchQuery(db *sql.DB, limit int, offset int, sortBy string, order string, searchQuery string) ([]models.Topic, error) {
	var topics []models.Topic
	args := []interface{}{searchQuery}

	query := `
	SELECT id, title, description, created_by, created_at
	FROM topics, plainto_tsquery('english', $1) AS query
	WHERE document @@ query
	`

	if sortBy == "relevance" {
		sortBy = "ts_rank(document, query)"
	}

	query = query + " ORDER BY " + sortBy + " " + order + " LIMIT $2 OFFSET $3"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)

	if err != nil {
		return topics, err
	}

	defer rows.Close()

	for rows.Next() {
		var topic models.Topic

		if err := rows.Scan(&topic.ID, &topic.Title, &topic.Description, &topic.CreatedBy, &topic.CreatedAt); err != nil {
			return topics, err
		}

		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return topics, err
	}

	return topics, nil
}

func ReadTopic(db *sql.DB, limit int, offset int, sortBy string, order string) ([]models.Topic, error) {
	var topics []models.Topic
	args := []interface{}{}

	query := `
	SELECT id, title, description, created_by, created_at
	FROM topics
	`

	query = query + " ORDER BY " + sortBy + " " + order + " LIMIT $1 OFFSET $2"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)

	if err != nil {
		return topics, err
	}

	defer rows.Close()

	for rows.Next() {
		var topic models.Topic

		if err := rows.Scan(&topic.ID, &topic.Title, &topic.Description, &topic.CreatedBy, &topic.CreatedAt); err != nil {
			return topics, err
		}

		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return topics, err
	}

	return topics, nil
}
