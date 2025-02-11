package routes

import (
	"gofast/handlers"
	"gofast/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome to GoFast API v1"})
		})

		/**
		 * Users routes
		 */
		api.POST("/users", services.CreateUser)
		api.GET("/users", services.GetUsers)
		api.GET("/users/:id", services.GetUserByID)
		api.PUT("/users/:id", services.UpdateUser)
		api.DELETE("/users/:id", services.DeleteUser)

		/**
		 * Assets routes
		 */
		api.POST("/assets", services.CreateAsset)
		api.GET("/assets", services.GetAssets)
		api.GET("/assets/:id", services.GetAssetByID)
		api.PUT("/assets/:id", services.UpdateAsset)
		api.DELETE("/assets/:id", services.DeleteAsset)
	}
}
