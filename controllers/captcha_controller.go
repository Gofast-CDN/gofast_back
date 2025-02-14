package controllers

import (
	"gofast/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaptchaController struct{}

func NewCaptchaController() *CaptchaController {
	return &CaptchaController{}
}

func (cc *CaptchaController) VerifyCaptcha(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide"})
		return
	}

	// Vérification du token reCAPTCHA
	valid, err := services.VerifyRecaptcha(req.Token)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Échec de la vérification reCAPTCHA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
