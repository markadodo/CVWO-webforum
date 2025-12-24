package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"

	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTopicHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateTopicInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Title == "" || input.Description == "" || input.CreatedBy <= 0 {
			c.JSON(400, gin.H{"error": "empty fields"})
			return
		}

		topic := models.Topic{
			Title:       input.Title,
			Description: input.Description,
			CreatedBy:   input.CreatedBy,
		}

		if err := database.CreateTopic(db, &topic); err != nil {
			c.JSON(500, gin.H{"error": "Could not create topic"})
			return
		}

		c.JSON(201, gin.H{
			"id":          topic.ID,
			"title":       topic.Title,
			"description": topic.Description,
			"created_by":  topic.CreatedBy,
			"created_at":  topic.CreatedAt,
		})
	}
}

func ReadTopicByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		topic, err := database.ReadTopicByID(db, id)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if topic == nil {
			c.JSON(404, gin.H{"error": "Topic not found"})
			return
		}

		c.JSON(200, gin.H{
			"id":          topic.ID,
			"title":       topic.Title,
			"description": topic.Description,
			"created_by":  topic.CreatedBy,
			"created_at":  topic.CreatedAt,
		})
	}
}

func UpdateTopicByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)
		if err != nil {
			c.JSON(404, gin.H{"error": "Invalid ID"})
			return
		}

		var input models.UpdateTopicInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Title != nil && *input.Title == "" {
			c.JSON(400, gin.H{"error": "Title cannot be empty"})
			return
		}
		if input.Description != nil && *input.Description == "" {
			c.JSON(400, gin.H{"error": "Description cannot be empty"})
			return
		}

		empty_update, topic_not_found, err := database.UpdateTopicByID(db, id, &input)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not update topic"})
			return
		}

		if empty_update {
			c.JSON(400, gin.H{"error": "Empty update"})
			return
		}

		if topic_not_found {
			c.JSON(404, gin.H{"error": "Topic not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Updated successfully"})
	}
}

func DeleteTopicByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		topic_not_found, err := database.DeleteTopicByID(db, id)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not delete topic"})
			return
		}

		if topic_not_found {
			c.JSON(404, gin.H{"error": "Topic not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Topic deleted"})
	}
}
