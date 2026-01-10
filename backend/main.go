package main

// import (
// 	// "net/http"

// 	"github.com/gin-gonic/gin"

// 	"log"

// 	"fmt"

// 	"backend/database"
// 	"backend/handlers"
// 	"backend/middleware"

// 	_ "github.com/lib/pq"
// )

// func main() {
// 	db, err := database.ConnectDB()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer db.Close()

// 	fmt.Println("db working")

// 	err = database.InitDB(db)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("db created")

// 	test := gin.Default()

// 	test.POST("/login", handlers.LoginHandler(db))
// 	test.GET("/topics/:topic_id/posts", handlers.ReadPostByTopicIDHandler(db))
// 	test.GET("/posts/:post_id/comments", handlers.ReadCommentByPostIDHandler(db))
// 	test.POST("/users", handlers.CreateUserHandler(db))
// 	test.GET("/topics/:topic_id/posts/search", handlers.ReadPostBySearchQueryHandler(db))
// 	test.GET("/topics", handlers.ReadTopicHandler(db))
// 	test.GET("/posts", handlers.ReadPostHandler(db))

// 	r := test.Group("/logged_in")

// 	r.Use(middleware.JWTAuthorisation())
// 	{
// 		r.POST("/users", handlers.CreateUserHandler(db))

// 		r.GET("/users/:user_id", middleware.CheckOwnershipByID(db, database.GetUserOwnerByID), handlers.ReadUserByIDHandler(db))

// 		r.PATCH("/users/:user_id", handlers.UpdateUserByIDHandler(db))

// 		r.DELETE("/users/:user_id", handlers.DeleteUserByIDHandler(db))

// 		r.POST("/topics", handlers.CreateTopicHandler(db))

// 		r.GET("/topics/:topic_id", middleware.CheckOwnershipByID(db, database.GetTopicOwnerByID), handlers.ReadTopicByIDHandler(db))

// 		r.PATCH("/topics/:topic_id", handlers.UpdateTopicByIDHandler(db))

// 		r.DELETE("/topics/:topic_id", handlers.DeleteTopicByIDHandler(db))

// 		r.POST("/posts", handlers.CreatePostHandler(db))

// 		r.GET("/posts/:post_id", handlers.ReadPostByIDHandler(db))

// 		r.PATCH("/posts/:post_id", handlers.UpdatePostByIDHandler(db))

// 		r.DELETE("/posts/:post_id", handlers.DeletePostByIDHandler(db))

// 		r.POST("/comments", handlers.CreateCommentHandler(db))

// 		r.GET("/comments/:comment_id", handlers.ReadCommentByIDHandler(db))

// 		r.PATCH("/comments/:comment_id", handlers.UpdateCommentByIDHandler(db))

// 		r.DELETE("/comments/:comment_id", handlers.DeleteCommentByIDHandler(db))

// 		r.POST("/posts/:post_id/reaction", handlers.CreatePostReactionHandler(db))

// 		r.DELETE("/posts/:post_id/reaction", handlers.DeletePostReactionHandler(db))

// 		r.POST("/comments/:comment_id/reaction", handlers.CreateCommentReactionHandler(db))

// 		r.DELETE("/comments/:comment_id/reaction", handlers.DeleteCommentReactionHandler(db))
// 	}

// 	test.Run(":8080")
// }

import (
	"github.com/gin-gonic/gin"

	"log"

	"fmt"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func main() {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	db, err := sql.Open("postgres", connStr)

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
	test.GET("/topics/:topic_id/posts", handlers.ReadPostByTopicIDHandler(db))
	test.GET("/posts/:post_id/comments", handlers.ReadCommentByPostIDHandler(db))
	test.POST("/users", handlers.CreateUserHandler(db))
	test.GET("/topics/:topic_id/posts/search", handlers.ReadPostBySearchQueryHandler(db))
	test.GET("/topics", handlers.ReadTopicHandler(db))
	test.GET("/posts", handlers.ReadPostHandler(db))

	r := test.Group("/logged_in")

	r.Use(middleware.JWTAuthorisation())
	{
		r.POST("/users", handlers.CreateUserHandler(db))

		r.GET("/users/:user_id", middleware.CheckOwnershipByID(db, database.GetUserOwnerByID), handlers.ReadUserByIDHandler(db))

		r.PATCH("/users/:user_id", handlers.UpdateUserByIDHandler(db))

		r.DELETE("/users/:user_id", handlers.DeleteUserByIDHandler(db))

		r.POST("/topics", handlers.CreateTopicHandler(db))

		r.GET("/topics/:topic_id", middleware.CheckOwnershipByID(db, database.GetTopicOwnerByID), handlers.ReadTopicByIDHandler(db))

		r.PATCH("/topics/:topic_id", handlers.UpdateTopicByIDHandler(db))

		r.DELETE("/topics/:topic_id", handlers.DeleteTopicByIDHandler(db))

		r.POST("/posts", handlers.CreatePostHandler(db))

		r.GET("/posts/:post_id", handlers.ReadPostByIDHandler(db))

		r.PATCH("/posts/:post_id", handlers.UpdatePostByIDHandler(db))

		r.DELETE("/posts/:post_id", handlers.DeletePostByIDHandler(db))

		r.POST("/comments", handlers.CreateCommentHandler(db))

		r.GET("/comments/:comment_id", handlers.ReadCommentByIDHandler(db))

		r.PATCH("/comments/:comment_id", handlers.UpdateCommentByIDHandler(db))

		r.DELETE("/comments/:comment_id", handlers.DeleteCommentByIDHandler(db))

		r.POST("/posts/:post_id/reaction", handlers.CreatePostReactionHandler(db))

		r.DELETE("/posts/:post_id/reaction", handlers.DeletePostReactionHandler(db))

		r.POST("/comments/:comment_id/reaction", handlers.CreateCommentReactionHandler(db))

		r.DELETE("/comments/:comment_id/reaction", handlers.DeleteCommentReactionHandler(db))
	}

	test.Run(":8080")
}
