package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"log"

	"fmt"

	"backend/database"
	"backend/handlers"
	"backend/middleware"

	_ "github.com/lib/pq"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to Postgres")

	err = database.InitDB(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB created")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local testing
	}

	router := gin.Default()

	routes := router.Group("/")
	routes.Use(middleware.EnableCORS())

	// Catching OPTIONS
	routes.OPTIONS("/*path")

	// PUBLIC ROUTES (No Authentication Required)
	public := routes.Group("/public")
	public.Use(middleware.JWTAuthorisationPublic())
	{
		//Return User ID
		public.GET("/auth/loginStatus", handlers.ReadLoggedInUserID(db))

		// Authentication Routes
		public.POST("/auth/register", handlers.CreateUserHandler(db))
		public.POST("/auth/login", handlers.LoginHandler(db))
		public.POST("/auth/logout", handlers.LogoutHandler(db))

		// User Routes - Read Only
		public.GET("/users/:user_id", handlers.ReadUsernameByIDHandler(db))

		// Topic Routes - Read Only
		public.GET("/topics", handlers.ReadTopicHandler(db))
		public.GET("/topics/:topic_id", handlers.ReadTopicByIDHandler(db))
		public.GET("/topics/search", handlers.ReadTopicBySearchQueryHandler(db))

		// Post Routes - Read Only (Public Feed)
		public.GET("/posts", handlers.ReadPostHandler(db))
		public.GET("/posts/:post_id", handlers.ReadPostByIDHandler(db))
		public.PATCH("/posts/:post_id", handlers.UpdatePostViewsByIDHandler(db))
		public.GET("/topics/:topic_id/posts", handlers.ReadPostByTopicIDHandler(db))
		public.GET("/topics/:topic_id/posts/search", handlers.ReadPostBySearchQueryHandler(db))

		// Comment Routes - Read Only
		public.GET("/posts/:post_id/comments", handlers.ReadCommentByPostIDHandler(db))
		public.GET("/comments/:parent_comment_id", handlers.ReadCommentByParentCommentIDHandler(db))
	}

	// PROTECTED ROUTES (Authentication Required)
	protected := routes.Group("/logged_in")
	protected.Use(middleware.JWTAuthorisation())
	{
		// USER CRUD
		protected.GET("/users/:user_id", middleware.CheckOwnershipByID(db, database.GetUserOwnerByID), handlers.ReadUserByIDHandler(db))
		protected.PATCH("/users/:user_id", middleware.CheckOwnershipByID(db, database.GetUserOwnerByID), handlers.UpdateUserByIDHandler(db))
		protected.DELETE("/users/:user_id", middleware.CheckOwnershipByID(db, database.GetUserOwnerByID), handlers.DeleteUserByIDHandler(db))

		//TOPIC CRUD
		protected.POST("/topics", handlers.CreateTopicHandler(db))
		protected.PATCH("/topics/:topic_id", middleware.CheckOwnershipByID(db, database.GetTopicOwnerByID), handlers.UpdateTopicByIDHandler(db))
		protected.DELETE("/topics/:topic_id", middleware.CheckOwnershipByID(db, database.GetTopicOwnerByID), handlers.DeleteTopicByIDHandler(db))

		//POST CRUD
		protected.POST("topics/:topic_id/posts", handlers.CreatePostHandler(db))
		protected.PATCH("/posts/:post_id", middleware.CheckOwnershipByID(db, database.GetPostOwnerByID), handlers.UpdatePostByIDHandler(db))
		protected.DELETE("/posts/:post_id", middleware.CheckOwnershipByID(db, database.GetPostOwnerByID), handlers.DeletePostByIDHandler(db))

		//COMMENT CRUD
		protected.POST("/comments", handlers.CreateCommentHandler(db))
		protected.PATCH("/comments/:comment_id", middleware.CheckOwnershipByID(db, database.GetCommentOwnerByID), handlers.UpdateCommentByIDHandler(db))
		protected.DELETE("/comments/:comment_id", middleware.CheckOwnershipByID(db, database.GetCommentOwnerByID), handlers.DeleteCommentByIDHandler(db))

		//POST REACTIONS
		protected.POST("/posts/:post_id/reactions", handlers.CreatePostReactionHandler(db))
		protected.DELETE("/posts/:post_id/reactions", handlers.DeletePostReactionHandler(db))
		protected.GET("/posts/:post_id/reactions", handlers.ReadPostReactionHandler(db))

		//COMMENT REACTIONS
		protected.POST("/comments/:comment_id/reactions", handlers.CreateCommentReactionHandler(db))
		protected.DELETE("/comments/:comment_id/reactions", handlers.DeleteCommentReactionHandler(db))
		protected.GET("/comments/:comment_id/reactions", handlers.ReadCommentReactionHandler(db))
	}

	router.Run(":" + port)
}
