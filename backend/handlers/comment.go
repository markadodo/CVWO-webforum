package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"errors"

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
		strid := c.Param("comment_id")
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
		strid := c.Param("comment_id")
		id, err := strconv.ParseInt(strid, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
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
		strid := c.Param("comment_id")
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

func ReadCommentByPostIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("post_id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil || postID <= 0 {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

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

		if sortBy != "created_at" && sortBy != "likes" {
			sortBy = "created_at"
		}

		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}

		commentsData, err := database.ReadCommentByPostID(db, postID, limit, offset, sortBy, order)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if len(commentsData) == 0 {
			c.JSON(200, gin.H{
				"count":    0,
				"page":     page,
				"limit":    limit,
				"sort_by":  sortBy,
				"order":    order,
				"comments": []models.Comment{},
			})
			return
		}

		c.JSON(200, gin.H{
			"count":    len(commentsData),
			"page":     page,
			"limit":    limit,
			"sort_by":  sortBy,
			"order":    order,
			"comments": commentsData,
		})
	}
}

func CreateCommentReactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("comment_id")
		commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid comment id"})
			return
		}

		if commentID <= 0 {
			c.JSON(400, gin.H{"error": "Invalid comment ID"})
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

		var input models.CreateCommentReactionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid payload"})
			return
		}

		commentReaction := models.CommentReaction{
			CommentID: commentID,
			UserID:    userID,
			Reaction:  input.Reaction,
		}

		err = database.CreateCommentReaction(db, &commentReaction)

		if err != nil {
			if errors.Is(err, database.ErrDuplicateCommentReaction) {
				c.JSON(409, gin.H{
					"error": "User has already reacted to this comment"})
				return
			}
			c.JSON(500, gin.H{"error": "Could not create reaction"})
			return
		}

		c.JSON(200, gin.H{"status": "Reaction created"})
	}
}

func DeleteCommentReactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("comment_id")
		commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid comment id"})
			return
		}

		if commentID <= 0 {
			c.JSON(400, gin.H{"error": "Invalid comment ID"})
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

		comment_reaction_not_found, err := database.DeleteCommentReactionByCommentIDAndUserID(db, commentID, userID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not delete reaction"})
			return
		}

		if comment_reaction_not_found {
			c.JSON(404, gin.H{"error": "Reaction not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Reaction deleted"})
	}
}
