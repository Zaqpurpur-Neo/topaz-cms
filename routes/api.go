package routes

// basically like laravel structure

import (
	"art-cms/controllers"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(g *gin.Engine) {
	api := g.Group("/api")
	{
		postController := controllers.NewPostController()
		api.POST("/posts", postController.CreatePost)
	}
}
