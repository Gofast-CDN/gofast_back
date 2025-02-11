package routes

import (
	"gofast/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/health", handlers.HealthCheck)
	r.GET("/mongodb-health", handlers.MongoDBHealthCheck) // Nouvelle route pour vérifier MongoDB

	api := r.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome to GoFast API v1"})
		})

		api.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello, World!"})
		})
	}
}
