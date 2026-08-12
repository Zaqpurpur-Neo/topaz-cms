package main

import (
	"embed"
	"fmt"
	"log"
	"topaz-workspace/config"
	"topaz-workspace/routes"

	"github.com/gin-gonic/gin"
)

var frontendAssets embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	routes.ApiRoutes(r)

	serverURL := "http://localhost" + config.PORT
	fmt.Printf("[INFO]: Dashboard running at %s\n", serverURL)
	log.Fatal(r.Run(config.PORT))
}
