package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"

	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreatePostInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Title == "" || input.Description == "" || input.TopicID <= 0 || input.CreatedBy <= 0 {
			c.JSON(400, gin.H{"error": "empty fields"})
			return
		}

		post := models.Post{
			Title:       input.Title,
			Description: input.Description,
			TopicID:     input.TopicID,
			CreatedBy:   input.CreatedBy,
		}

		if err := database.CreatePost(db, &post); err != nil {
			c.JSON(500, gin.H{"error": "Could not create post"})
			return
		}

		c.JSON(201, gin.H{
			"id":          post.ID,
			"title":       post.Title,
			"description": post.Description,
			"topic_id":    post.TopicID,
			"likes":       post.Likes,
			"dislikes":    post.Dislikes,
			"is_edited":   post.IsEdited,
			"views":       post.Views,
			"popularity":  post.Popularity,
			"created_by":  post.CreatedBy,
			"created_at":  post.CreatedAt,
		})
	}
}

func ReadPostByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		post, err := database.ReadPostByID(db, id)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if post == nil {
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		}

		c.JSON(201, gin.H{
			"id":          post.ID,
			"title":       post.Title,
			"description": post.Description,
			"topic_id":    post.TopicID,
			"likes":       post.Likes,
			"dislikes":    post.Dislikes,
			"is_edited":   post.IsEdited,
			"views":       post.Views,
			"popularity":  post.Popularity,
			"created_by":  post.CreatedBy,
			"created_at":  post.CreatedAt,
		})
	}
}

func UpdatePostByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)
		if err != nil {
			c.JSON(404, gin.H{"error": "Invalid ID"})
			return
		}

		var input models.UpdatePostInput

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

		empty_update, post_not_found, err := database.UpdatePostByID(db, id, &input)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not update post"})
			return
		}

		if empty_update {
			c.JSON(400, gin.H{"error": "Empty update"})
			return
		}

		if post_not_found {
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Updated successfully"})
	}
}

func DeletePostByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		post_not_found, err := database.DeletePostByID(db, id)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not delete post"})
			return
		}

		if post_not_found {
			c.JSON(404, gin.H{"error": "Post not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Post deleted"})
	}
}
