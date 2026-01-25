package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"strings"

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

		if input.Title == "" || input.Description == "" {
			c.JSON(400, gin.H{"error": "empty fields"})
			return
		}
		userIDVal, exists := c.Get("user_id")

		if !exists {
			c.JSON(401, gin.H{"error": "Not logged in"})
			return
		}

		userID, match := userIDVal.(int64)

		if !match {
			c.JSON(401, gin.H{"error": "Invalid user ID"})
			return
		}

		topic := models.Topic{
			Title:       input.Title,
			Description: input.Description,
			CreatedBy:   userID,
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
		strid := c.Param("topic_id")
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
		strid := c.Param("topic_id")
		id, err := strconv.ParseInt(strid, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
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
		strid := c.Param("topic_id")
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

func ReadTopicBySearchQueryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid page"})
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid limit"})
			return
		}

		if page <= 0 {
			page = 1
		}

		if limit < 10 || limit >= 100 {
			limit = 10
		}

		offset := (page - 1) * limit

		sortBy := c.DefaultQuery("sort_by", "created_at")
		order := c.DefaultQuery("order", "DESC")

		if sortBy != "created_at" && sortBy != "relevance" {
			sortBy = "created_at"
		}

		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}

		searchQuery := c.DefaultQuery("q", "")
		searchQuery = strings.TrimSpace(searchQuery)
		if searchQuery == "" {
			c.JSON(400, gin.H{"error": "Query cannot be empty"})
			return
		}

		topicsData, err := database.ReadTopicBySearchQuery(db, limit, offset, sortBy, order, searchQuery)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if len(topicsData) == 0 {
			c.JSON(200, gin.H{
				"count":        0,
				"page":         page,
				"limit":        limit,
				"sort_by":      sortBy,
				"order":        order,
				"search_query": searchQuery,
				"topics":       []models.Topic{},
			})
			return
		}

		c.JSON(200, gin.H{
			"count":        len(topicsData),
			"page":         page,
			"limit":        limit,
			"sort_by":      sortBy,
			"order":        order,
			"search_query": searchQuery,
			"topics":       topicsData,
		})
	}
}

func ReadTopicHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid page"})
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid limit"})
			return
		}

		if page <= 0 {
			page = 1
		}

		if limit < 10 || limit >= 100 {
			limit = 10
		}

		offset := (page - 1) * limit

		sortBy := c.DefaultQuery("sort_by", "created_at")
		order := c.DefaultQuery("order", "DESC")

		if sortBy != "created_at" {
			sortBy = "created_at"
		}

		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}

		topicsData, err := database.ReadTopic(db, limit, offset, sortBy, order)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if len(topicsData) == 0 {
			c.JSON(200, gin.H{
				"count":   0,
				"page":    page,
				"limit":   limit,
				"sort_by": sortBy,
				"order":   order,
				"topics":  []models.Topic{},
			})
			return
		}

		c.JSON(200, gin.H{
			"count":   len(topicsData),
			"page":    page,
			"limit":   limit,
			"sort_by": sortBy,
			"order":   order,
			"topics":  topicsData,
		})
	}
}
