package main

import (
	"log"

	"gofast/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("server startup failed:", err)
	}
}
