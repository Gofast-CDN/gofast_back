package routes

import (
	"gofast/config"
	"gofast/controllers"
	"gofast/handlers"
	"gofast/middleware"
	"gofast/services"

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
}

// Global routes API
func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to GoFast API v1"})
	})

	setupUserRoutes(api)
	// Add other routes here
	// Ajouter la route de validation reCAPTCHA ici
	setupCaptchaRoutes(api)
}

// Ajouter une route pour le reCAPTCHA
func setupCaptchaRoutes(rg *gin.RouterGroup) {
	captcha := rg.Group("/captcha")
	{
		// Route pour valider le reCAPTCHA
		captcha.POST("/verify-recaptcha", verifyCaptcha)
	}
}

func verifyCaptcha(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Vérification du token reCAPTCHA
	valid, err := services.VerifyRecaptcha(req.Token)
	if err != nil || !valid {
		c.JSON(401, gin.H{"success": false, "message": "reCAPTCHA verification failed"})
		return
	}

	c.JSON(200, gin.H{"success": true})
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
	{
		assets.POST("", assetsController.CreateAsset)
		assets.GET("", assetsController.GetAssets)
		assets.GET("/:id", assetsController.GetAssetByID)
		assets.PUT("/:id", assetsController.UpdateAsset)
		assets.DELETE("/:id", assetsController.DeleteAsset)
	}
}
