package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"topaz-workspace/config"
	"topaz-workspace/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
		if filepath.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rute API tidak ditemukan"})
			return
		}
		c.File(filepath.Join("frontend/dist", "index.html"))
	})

	routes.ApiRoutes(r)

	serverURL := "http://localhost" + config.PORT
	fmt.Printf("[INFO]: Topaz Gin Server running at %s\n", serverURL)

	log.Fatal(r.Run(config.PORT))
}
