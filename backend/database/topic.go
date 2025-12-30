package database

import (
	"backend/models"
	"database/sql"
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
	VALUES (?, ?, ?, ?);
	`
	result, err := db.Exec(
		query,
		topic.Title,
		topic.Description,
		topic.CreatedBy,
		topic.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	topic.ID = id

	return nil
}

func ReadTopicByID(db *sql.DB, id int64) (*models.Topic, error) {
	topic := models.Topic{}

	query := `
	SELECT id, title, description, created_by, created_at
	FROM topics
	WHERE id = ?
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

	if input.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *input.Title)
	}

	if input.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *input.Description)
	}

	if len(updates) == 0 {
		return true, false, nil
	}

	query := "UPDATE topics SET " + strings.Join(updates, ", ") + " WHERE id = ?"
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
	query := "DELETE FROM topics WHERE id = ?"
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
