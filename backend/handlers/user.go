package handlers

import (
	"backend/database"
	"backend/models"
	"database/sql"

	"github.com/gin-gonic/gin"
)

//learning purpose
/*
func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Allow only POST
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode JSON body
		var input models.CreateUserInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "issue", http.StatusBadRequest)
			log.Println(err)
			return
		}

		// Basic validation
		if input.Username == "" || input.Password == "" {
			http.Error(w, "username and password required", http.StatusBadRequest)
			return
		}

		// Convert input → domain model
		user := models.User{
			Username: input.Username,
			Password: input.Password,
		}

		// Save to DB
		if err := database.CreateUser(db, &user); err != nil {
			http.Error(w, "could not create user", http.StatusInternalServerError)
			return
		}

		// Prepare response
		response := map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
*/

func CreateUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateUserInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if input.Password == "" || input.Username == "" {
			c.JSON(400, gin.H{"error": "Username or Password required"})
		}

		user := models.User{
			Password: input.Password,
			Username: input.Username,
		}

		if err := database.CreateUser(db, &user); err != nil {
			c.JSON(500, gin.H{"error": "Could not create user"})
		}

		c.JSON(201, gin.H{
			"username": user.Username,
		})
	}
}
