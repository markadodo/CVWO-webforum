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
	test.GET("/topics/:topic_id/posts", handlers.ReadPostByTopicIDHandler(db))
	test.GET("/posts/:post_id/comments", handlers.ReadCommentByPostIDHandler(db))
	test.POST("/users", handlers.CreateUserHandler(db))
	test.GET("/topics/:topic_id/posts/search", handlers.ReadPostBySearchQueryHandler(db))

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

// package main

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	db, _ := sql.Open("sqlite3", ":memory:")
// 	defer db.Close()

// 	_, err := db.Exec("CREATE VIRTUAL TABLE test USING fts5(content);")
// 	if err != nil {
// 		fmt.Println("FTS5 NOT enabled:", err)
// 	} else {
// 		fmt.Println("FTS5 IS enabled!")
// 	}
// }
