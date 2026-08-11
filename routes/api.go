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

		kritaRoutePath := "/krita"
		kritaWorkspaceController := controllers.NewKritaWorkspaceController(kritaRoutePath, "/api")
		api.GET(kritaRoutePath, kritaWorkspaceController.GetKritaWorkspace)
		api.GET(kritaRoutePath+"/:workspace", kritaWorkspaceController.GetKritaWorkspaceByName)
		api.GET(kritaRoutePath+"/:workspace/preview/:artwork", kritaWorkspaceController.GetArtworkPreview)
	}
}
