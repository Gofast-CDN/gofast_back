package controllers

import (
	"gofast/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaptchaController struct {
	captchaService *services.CaptchaService
}

// Constructeur du contrôleur
func NewCaptchaController() *CaptchaController {
	return &CaptchaController{
		captchaService: services.NewCaptchaService(),
	}
}

// Vérifie un token reCAPTCHA envoyé depuis le client
func (cc *CaptchaController) VerifyCaptcha(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	// Validation de la requête
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide", "details": err.Error()})
		return
	}

	// Vérification du token via le service reCAPTCHA
	valid, err := cc.captchaService.VerifyRecaptcha(req.Token)
	if err != nil || !valid { // ✅ Utilisation de `valid`
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Échec de la vérification reCAPTCHA", "error": err.Error()})
		return
	}

	// Réponse en cas de succès
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "reCAPTCHA validé avec succès"})
}
