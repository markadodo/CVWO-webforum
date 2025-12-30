package middleware

import (
	"backend/auth"
	"database/sql"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthorisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		token, err := auth.CheckTokenValidity(tokenStr)

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, match := token.Claims.(jwt.MapClaims); match {
			if userID, match := claims["user_id"].(float64); match {
				c.Set("user_id", int64(userID))

			} else {
				c.JSON(401, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
		} else {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

type resourceFetcher func(*sql.DB, int64) (int64, error)

func CheckOwnershipByID(db *sql.DB, fetcher resourceFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceIDStr := c.Param("id")

		resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			c.Abort()
			return
		}

		ownerUserID, err := fetcher(db, resourceID)

		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if ownerUserID == 0 {
			c.JSON(404, gin.H{"error": "Resource not found"})
			c.Abort()
			return
		}

		currentUserID, exists := c.Get("user_id")

		if !exists {
			c.JSON(401, gin.H{"error": "Not logged in"})
			c.Abort()
			return
		}

		currentUserID, match := currentUserID.(int64)

		if !match {
			c.JSON(401, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		if currentUserID != ownerUserID {
			c.JSON(403, gin.H{"error": "Unauthorised"})
			c.Abort()
			return
		}

		c.Next()
	}
}
