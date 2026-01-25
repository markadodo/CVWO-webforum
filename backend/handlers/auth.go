package handlers

import (
	"database/sql"

	"os"

	"github.com/gin-gonic/gin"

	"backend/auth"
	"backend/database"
	"backend/models"
)

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var loginData models.LoginUserData

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		userData, err := database.ReadUserByUsername(db, loginData.Username)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if userData == nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		userID, err := auth.CheckLoginValidity(userData, &loginData)

		if err != nil {
			c.JSON(400, gin.H{"error": "Wrong password"})
			return
		}

		tokenStr, err := auth.GenerateJWT(userID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not generate token"})
			return
		}

		var current_status bool = false

		if status := os.Getenv("STATUS"); status == "deployment" {
			current_status = true
		}

		c.SetCookie(
			"token",
			tokenStr,
			3600,
			"/",
			"",
			current_status,
			true,
		)

		c.JSON(200, gin.H{"user_id": userID})
	}
}

func LogoutHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.JSON(200, gin.H{"status": "Logged out"})
	}
}

func ReadLoggedInUserID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if userIDval, exists := c.Get("user_id"); exists {
			userID := userIDval.(int64)
			user, err := database.ReadUserByID(db, userID)

			if err != nil {
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			}

			if user == nil {
				c.JSON(404, gin.H{"error": "User not found"})
				return
			}
			c.JSON(200, gin.H{
				"logged_in": true,
				"user_id":   userID,
				"username":  user.Username,
			})
		} else {
			c.JSON(400, gin.H{
				"logged_in": false,
				"user_id":   -1,
				"username":  "",
			})
		}
	}
}
