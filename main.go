package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"topaz-workspace/config"
	"topaz-workspace/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var frontendStatic embed.FS

// i need to add frontend assets, from the frontend/dist folder
func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	r.StaticFile("/topaz.webp", "./frontend/dist/topaz.webp")

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:7733", "http://localhost:5173", "http://localhost:4173"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.NoRoute(func(c *gin.Context) {
		// Pastikan rute yang salah ketik di API tetap mengembalikan error 404 API,
		// bukan malah mengirim file HTML frontend
		if filepath.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rute API tidak ditemukan"})
			return
		}

		// Kirim index.html untuk diurus oleh React/Vue/Svelte router
		c.File(filepath.Join("frontend/dist", "index.html"))
	})

	routes.ApiRoutes(r)

	serverURL := "http://localhost" + config.PORT
	fmt.Printf("[INFO]: Dashboard running at %s\n", serverURL)
	log.Fatal(r.Run(config.PORT))
}
