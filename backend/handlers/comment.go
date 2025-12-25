package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"

	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCommentHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateCommentInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Description == "" || input.PostID <= 0 || (input.ParentCommentID != nil && *input.ParentCommentID <= 0) || input.CreatedBy <= 0 {
			c.JSON(400, gin.H{"error": "empty fields"})
			return
		}

		comment := models.Comment{
			Description:     input.Description,
			PostID:          input.PostID,
			ParentCommentID: input.ParentCommentID,
			CreatedBy:       input.CreatedBy,
		}

		if err := database.CreateComment(db, &comment); err != nil {
			c.JSON(500, gin.H{"error": "Could not create comment"})
			return
		}

		var checked_parent_comment_id interface{}

		if comment.ParentCommentID == nil {
			checked_parent_comment_id = nil
		} else {
			checked_parent_comment_id = *comment.ParentCommentID
		}

		c.JSON(201, gin.H{
			"id":                comment.ID,
			"description":       comment.Description,
			"likes":             comment.Likes,
			"dislikes":          comment.Dislikes,
			"is_edited":         comment.IsEdited,
			"post_id":           comment.PostID,
			"parent_comment_id": checked_parent_comment_id,
			"created_by":        comment.CreatedBy,
			"created_at":        comment.CreatedAt,
		})
	}
}

func ReadCommentByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		comment, err := database.ReadCommentByID(db, id)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if comment == nil {
			c.JSON(404, gin.H{"error": "Comment not found"})
			return
		}

		var checked_parent_comment_id interface{}

		if comment.ParentCommentID == nil {
			checked_parent_comment_id = nil
		} else {
			checked_parent_comment_id = *comment.ParentCommentID
		}

		c.JSON(200, gin.H{
			"id":                comment.ID,
			"description":       comment.Description,
			"likes":             comment.Likes,
			"dislikes":          comment.Dislikes,
			"is_edited":         comment.IsEdited,
			"post_id":           comment.PostID,
			"parent_comment_id": checked_parent_comment_id,
			"created_by":        comment.CreatedBy,
			"created_at":        comment.CreatedAt,
		})
	}
}

func UpdateCommentByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)
		if err != nil {
			c.JSON(404, gin.H{"error": "Invalid ID"})
			return
		}

		var input models.UpdateCommentInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Description != nil && *input.Description == "" {
			c.JSON(400, gin.H{"error": "Description cannot be empty"})
			return
		}

		empty_update, comment_not_found, err := database.UpdateCommentByID(db, id, &input)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not update comment"})
			return
		}

		if empty_update {
			c.JSON(400, gin.H{"error": "Empty update"})
			return
		}

		if comment_not_found {
			c.JSON(404, gin.H{"error": "Comment not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Updated successfully"})
	}
}

func DeleteCommentByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		comment_not_found, err := database.DeleteCommentByID(db, id)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not delete comment"})
			return
		}

		if comment_not_found {
			c.JSON(404, gin.H{"error": "Comment not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Comment deleted"})
	}
}
