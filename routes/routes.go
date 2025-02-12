package routes

import (
	"gofast/config"
	"gofast/controllers"
	"gofast/handlers"
	"gofast/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	r.Use(middleware.CorsMiddleware())
	setupHealthRoutes(r)
	setupAPIRoutes(r)
}

func setupHealthRoutes(r *gin.Engine) {
	r.GET("/", handlers.HealthCheck)
}

// Global routes API
func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to GoFast API v1"})
	})

	setupUserRoutes(api)
	// Add other routes here
}

func setupUserRoutes(rg *gin.RouterGroup) {
	userController := controllers.NewUserController()
	users := rg.Group("/users")
	{
		// Public routes
		users.POST("/register", userController.Register)
		users.POST("/login", userController.Login)

		// Protected routes
		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", userController.GetMe)
		}
	}
}
