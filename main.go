package main

import (
	"art-cms/routes"
	"embed"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var frontendAssets embed.FS

const PORT = "7733"
const REPO_PATH = "../art-gallery-data"

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	routes.ApiRoutes(r)

	serverURL := "http://localhost:" + PORT
	fmt.Printf("[INFO]: Dashboard running at %s\n", serverURL)
	log.Fatal(r.Run(":" + PORT))
}
