package main

import (
	"log"

	"gofast/config"
	"gofast/database"
	"gofast/routes"
	"gofast/services"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}
	// Connexion à la base de données
	database.Connect()

	// Connexion au service Azure Blob Storage
	_, err := services.NewBlobStorageService()
	if err != nil {
		log.Fatalf("Failed to initialize blob storage service: %v", err)
	}

	// Initialisation du routeur Gin
	r := gin.Default()

	// Configuration des routes
	routes.SetupRoutes(r)

	// Démarrer le serveur
	if err := r.Run(":80"); err != nil {
		log.Fatal("server startup failed:", err)
	}
}
