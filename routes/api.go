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

		kritaArtworkController := controllers.NewKritaArtworkController(kritaRoutePath, "/api")
		api.GET(kritaRoutePath+"/:workspace/preview/:artwork", kritaArtworkController.GetArtworkPreview)
		api.GET(kritaRoutePath+"/:workspace/original/:artwork", kritaArtworkController.GetArtworkOriginal)
		api.GET(kritaRoutePath+"/:workspace/references/:artwork/:refname", kritaArtworkController.GetArtworkReferences)
		api.GET(kritaRoutePath+"/:workspace/board/:artwork", kritaArtworkController.GetBoard)
	}
}
