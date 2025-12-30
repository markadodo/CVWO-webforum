package main

import (
	// "net/http"

	"github.com/gin-gonic/gin"

	"log"

	"fmt"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
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

	test := gin.Default()

	test.POST("/login", handlers.LoginHandler(db))

	r := test.Group("/logged_in")

	r.Use(middleware.JWTAuthorisation())
	{
		r.POST("/users", handlers.CreateUserHandler(db))

		r.GET("/users/:id", middleware.CheckOwnershipByID(db, database.GetUserOwnerByID), handlers.ReadUserByIDHandler(db))

		r.PATCH("/users/:id", handlers.UpdateUserByIDHandler(db))

		r.DELETE("/users/:id", handlers.DeleteUserByIDHandler(db))

		r.POST("/topics", handlers.CreateTopicHandler(db))

		r.GET("/topics/:id", middleware.CheckOwnershipByID(db, database.GetTopicOwnerByID), handlers.ReadTopicByIDHandler(db))

		r.PATCH("/topics/:id", handlers.UpdateTopicByIDHandler(db))

		r.DELETE("/topics/:id", handlers.DeleteTopicByIDHandler(db))

		r.POST("/posts", handlers.CreatePostHandler(db))

		r.GET("/posts/:id", handlers.ReadPostByIDHandler(db))

		r.PATCH("/posts/:id", handlers.UpdatePostByIDHandler(db))

		r.DELETE("/posts/:id", handlers.DeletePostByIDHandler(db))

		r.POST("/comments", handlers.CreateCommentHandler(db))

		r.GET("/comments/:id", handlers.ReadCommentByIDHandler(db))

		r.PATCH("/comments/:id", handlers.UpdateCommentByIDHandler(db))

		r.DELETE("/comments/:id", handlers.DeleteCommentByIDHandler(db))
	}

	test.Run(":8080")
}
