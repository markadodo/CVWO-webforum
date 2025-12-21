package main

import (
	// "net/http"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"

	"fmt"

	"backend/database"
	"backend/handlers"
)

func main() {
	db, err := database.ConnectDB()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Println("db working")

	err = database.InitDB(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("db created")

	r := gin.Default()

	r.GET("/testing", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("/user", handlers.CreateUserHandler(db))

	r.Run(":8080")
}
