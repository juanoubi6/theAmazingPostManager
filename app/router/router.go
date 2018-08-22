package router

import (
	"github.com/aviddiviner/gin-limit"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingPostManager/app/config"
	"theAmazingPostManager/app/middleware"
	"theAmazingPostManager/app/controllers/post"
	"theAmazingPostManager/app/controllers/comment"
)

var router *gin.Engine

func CreateRouter() {
	router = gin.New()

	router.Use(gin.Logger())
	router.Use(nice.Recovery(recoveryHandler))
	router.Use(limit.MaxAllowed(10))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET,PUT,POST,DELETE"},
		AllowHeaders:    []string{"accept,x-access-token,content-type,authorization"},
	}))

	public := router.Group("/")
	{
		public.GET("/post", middleware.Paginate(), post.GetAllPosts)
		public.GET("/post/:id", post.GetPost)
		//Ultimos 10 post creados
	}

	postCreation := router.Group("/post", middleware.ValidateToken())
	{
		postCreation.POST("", post.CreatePost)
		postCreation.PUT("/:id", post.ModifyPost)
		postCreation.DELETE("/:id", post.DeletePost)
		//Votar post
		postCreation.PATCH("/:id", post.Vote)
	}

	commentCreation := router.Group("/post/:postID", middleware.ValidateToken())
	{
		commentCreation.POST("/", comment.AddComment)
		commentCreation.PUT("/:id", comment.EditComment)
		commentCreation.DELETE("/:id", comment.DeleteComment)
		//Rate comment
	}



}

func RunRouter() {
	router.Run(":" + config.GetConfig().PORT)
}

func recoveryHandler(c *gin.Context, err interface{}) {
	detail := ""
	if config.GetConfig().ENV == "develop" {
		detail = err.(error).Error()
	}
	c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": detail})
}
