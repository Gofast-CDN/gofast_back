package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type CaptchaService struct{}

// Constructeur du service
func NewCaptchaService() *CaptchaService {
	return &CaptchaService{}
}

// Vérifie le reCAPTCHA auprès de Google
func (cs *CaptchaService) VerifyRecaptcha(token string) (bool, error) {
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		return false, fmt.Errorf("clé secrète reCAPTCHA manquante")
	}

	// URL de l'API de vérification reCAPTCHA
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	// Préparation des données
	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", token)

	// Requête HTTP avec timeout
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la requête reCAPTCHA : %v", err)
	}
	defer resp.Body.Close()

	// Lecture de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la lecture de la réponse reCAPTCHA : %v", err)
	}

	// Parsing JSON
	var result struct {
		Success     bool     `json:"success"`
		ChallengeTs string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		ErrorCodes  []string `json:"error-codes"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("erreur lors du parsing JSON : %v", err)
	}

	// Vérification de la réponse
	if !result.Success {
		return false, fmt.Errorf("échec de la validation reCAPTCHA, erreurs: %v", result.ErrorCodes)
	}

	return true, nil
}
