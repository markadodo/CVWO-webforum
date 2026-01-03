package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"

	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateUserInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Username == "" || input.Password == "" {
			c.JSON(400, gin.H{"error": "empty fields"})
			return
		}

		user := models.User{
			Password: input.Password,
			Username: input.Username,
		}

		if err := database.CreateUser(db, &user); err != nil {
			c.JSON(500, gin.H{"error": "Could not create user"})
			return
		}

		c.JSON(201, gin.H{
			"id":       user.ID,
			"username": user.Username,
		})
	}
}

func ReadUserByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("user_id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		user, err2 := database.ReadUserByID(db, id)
		//consider emptying the passwordhash field
		if err2 != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if user == nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		c.JSON(200, gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"created_at":  user.CreatedAt,
			"last_active": user.LastActive,
		})
	}
}

func UpdateUserByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("user_id")
		id, err := strconv.ParseInt(strid, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		var input models.UpdateUserInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Username != nil && *input.Username == "" {
			c.JSON(400, gin.H{"error": "Username cannot be empty"})
			return
		}
		if input.Password != nil && *input.Password == "" {
			c.JSON(400, gin.H{"error": "Password cannot be empty"})
			return
		}

		empty_update, user_not_found, err := database.UpdateUserByID(db, id, &input)
		//consider emptying the password field

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not update user"})
			return
		}

		if empty_update {
			c.JSON(400, gin.H{"error": "Empty update"})
			return
		}

		if user_not_found {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		c.JSON(200, gin.H{"status": "Updated successfully"})
	}
}

func DeleteUserByIDHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		strid := c.Param("user_id")
		id, err := strconv.ParseInt(strid, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		user_not_found, err := database.DeleteUserByID(db, id)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not delete user"})
			return
		}

		if user_not_found {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		c.JSON(200, gin.H{"status": "User deleted"})
	}
}
