package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"errors"
	"strings"

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
		strid := c.Param("post_id")
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

		c.JSON(200, gin.H{
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
		strid := c.Param("post_id")
		id, err := strconv.ParseInt(strid, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
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
		strid := c.Param("post_id")
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

func ReadPostByTopicIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicIDStr := c.Param("topic_id")
		topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
		if err != nil || topicID <= 0 {
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

		if sortBy != "created_at" && sortBy != "popularity" && sortBy != "views" {
			sortBy = "created_at"
		}

		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}

		postsData, err := database.ReadPostByTopicID(db, topicID, limit, offset, sortBy, order)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if len(postsData) == 0 {
			c.JSON(200, gin.H{
				"count":   0,
				"page":    page,
				"limit":   limit,
				"sort_by": sortBy,
				"order":   order,
				"posts":   []models.Post{},
			})
			return
		}

		c.JSON(200, gin.H{
			"count":   len(postsData),
			"page":    page,
			"limit":   limit,
			"sort_by": sortBy,
			"order":   order,
			"posts":   postsData,
		})
	}
}

func ReadPostBySearchQueryHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicIDStr := c.Param("topic_id")
		var topicID int64 = 0
		if topicIDStr != "" {
			parsed, err := strconv.ParseInt(topicIDStr, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid ID"})
				return
			}
			topicID = parsed
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

		if sortBy != "created_at" && sortBy != "popularity" && sortBy != "views" && sortBy != "relevance" {
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

		postsData, err := database.ReadPostBySearchQuery(db, topicID, limit, offset, sortBy, order, searchQuery)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if len(postsData) == 0 {
			c.JSON(200, gin.H{
				"count":        0,
				"page":         page,
				"limit":        limit,
				"sort_by":      sortBy,
				"order":        order,
				"search_query": searchQuery,
				"posts":        []models.Post{},
			})
			return
		}

		c.JSON(200, gin.H{
			"count":        len(postsData),
			"page":         page,
			"limit":        limit,
			"sort_by":      sortBy,
			"order":        order,
			"search_query": searchQuery,
			"posts":        postsData,
		})
	}
}

func CreatePostReactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("post_id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid post id"})
			return
		}

		if postID <= 0 {
			c.JSON(400, gin.H{"error": "Invalid post ID"})
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

		var input models.CreatePostReactionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid payload"})
			return
		}

		postReaction := models.PostReaction{
			PostID:   postID,
			UserID:   userID,
			Reaction: input.Reaction,
		}

		err = database.CreatePostReaction(db, &postReaction)

		if err != nil {
			if errors.Is(err, database.ErrDuplicatePostReaction) {
				c.JSON(409, gin.H{"error": "User has already reacted to this post"})
				return
			}
			c.JSON(500, gin.H{"error": "Could not create reaction"})
			return
		}

		c.JSON(200, gin.H{"status": "Reaction created"})
	}
}

func DeletePostReactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("post_id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid post id"})
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

		post_reaction_not_found, err := database.DeletePostReactionByPostIDAndUserID(db, postID, userID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not delete reaction"})
			return
		}

		if post_reaction_not_found {
			c.JSON(404, gin.H{"error": "Reaction not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Reaction deleted"})
	}
}

func ReadPostHandler(db *sql.DB) gin.HandlerFunc {
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

		if sortBy != "created_at" && sortBy != "popularity" && sortBy != "views" {
			sortBy = "created_at"
		}

		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}

		postsData, err := database.ReadPost(db, limit, offset, sortBy, order)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if len(postsData) == 0 {
			c.JSON(200, gin.H{
				"count":   0,
				"page":    page,
				"limit":   limit,
				"sort_by": sortBy,
				"order":   order,
				"posts":   []models.Post{},
			})
			return
		}

		c.JSON(200, gin.H{
			"count":   len(postsData),
			"page":    page,
			"limit":   limit,
			"sort_by": sortBy,
			"order":   order,
			"posts":   postsData,
		})
	}
}
