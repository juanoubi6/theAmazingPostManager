package router

import (
	"github.com/aviddiviner/gin-limit"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingPostManager/app/config"
	"theAmazingPostManager/app/controllers/comment"
	"theAmazingPostManager/app/controllers/post"
	"theAmazingPostManager/app/middleware"
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
		public.GET("/posts/:postID/comment/:id", comment.GetComment)
		public.GET("/lastPosts", post.GetLastPosts)
		public.GET("/lastComments", comment.GetLastComments)
	}

	postCreation := router.Group("/post", middleware.ValidateToken())
	{
		postCreation.POST("", post.CreatePost)
		postCreation.PUT("/:id", middleware.IsPostOwner(false), post.ModifyPost)
		postCreation.DELETE("/:id", middleware.IsPostOwner(true), post.DeletePost)
		postCreation.PATCH("/:id", post.VotePost)
	}

	commentCreation := router.Group("/posts/:postID/comment", middleware.ValidateToken())
	{
		commentCreation.POST("", comment.AddComment)
		commentCreation.PUT("/:id", middleware.IsCommentOwner(false), comment.EditComment)
		commentCreation.DELETE("/:id", middleware.IsCommentOwner(true), comment.DeleteComment)
		commentCreation.PATCH("/:id", comment.VoteComment)
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
