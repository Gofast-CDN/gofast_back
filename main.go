package main

import (
	"log"

	"gofast/database"
	"gofast/routes"
	"gofast/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connexion à la base de données
	database.Connect()
	// Connexion au service Azure Blob Storage
	_, err := services.NewBlobStorageService()
	if err != nil {
		log.Fatalf("Failed to initialize blob storage service: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	if err := r.Run(":80"); err != nil {
		log.Fatal("server startup failed:", err)
	}
}
