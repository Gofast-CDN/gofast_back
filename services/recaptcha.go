package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RecaptchaResponse struct {
	Success     bool   `json:"success"`
	ChallengeTs string `json:"challenge_ts"`
	Hostname    string `json:"hostname"`
}

func VerifyRecaptcha(token string) (bool, error) {
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		return false, fmt.Errorf("clé secrète reCAPTCHA manquante")
	}

	// URL de l'API reCAPTCHA pour la vérification
	url := "https://www.google.com/recaptcha/api/siteverify"

	// Crée le corps de la requête
	data := fmt.Sprintf("secret=%s&response=%s", secretKey, token)

	// Effectue la requête HTTP vers l'API reCAPTCHA
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Effectue la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Lire la réponse JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result RecaptchaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	// Vérifie si la validation a réussi
	if result.Success {
		return true, nil
	}

	return false, fmt.Errorf("Échec de la validation reCAPTCHA")
}
