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
	r.GET("/health", handlers.HealthCheck)
	r.GET("/.well-known/acme-challenge", handlers.HealthCheck)
	r.GET("/.well-known/acme-challenge/:id", handlers.HealthCheck)
}

// Global routes API
func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to GoFast API v1"})
	})

	setupUserRoutes(api)
	setupCaptchaRoutes(api) // Route reCAPTCHA
	setupAssetsRoutes(api)
}

// Routes pour la validation du reCAPTCHA
func setupCaptchaRoutes(rg *gin.RouterGroup) {
	captchaController := controllers.NewCaptchaController()
	captcha := rg.Group("/captcha")
	{
		captcha.POST("/verify-recaptcha", captchaController.VerifyCaptcha)
	}
}

func setupUserRoutes(rg *gin.RouterGroup) {
	userController := controllers.NewUserController()
	users := rg.Group("/users")
	{
		// Routes publiques
		users.POST("/register", userController.Register)
		users.POST("/login", userController.Login)
		users.DELETE("/delete/:id", userController.Delete)

		// Routes protégées
		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", userController.GetMe)
		}
	}
}

func setupAssetsRoutes(rg *gin.RouterGroup) {
	assetsController := controllers.NewAssetsController()
	assets := rg.Group("/assets")
	assets.Use(middleware.AuthMiddleware())
	{
		assets.POST("", assetsController.CreateFileAsset)
		assets.POST("/folder", assetsController.CreateRepoAsset)
		assets.GET("/recent", assetsController.GetRecentAssetsFiles)
		assets.GET("/folder/recent", assetsController.GetRecentAssetsFolder)
		assets.GET("", assetsController.GetAssets)
		assets.GET("/:id", assetsController.GetAssetByID)
		assets.PUT("/:id", assetsController.UpdateAsset)
		assets.DELETE("/:id", assetsController.DeleteAsset)
	}
}
