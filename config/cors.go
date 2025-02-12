package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type CorsConfig struct {
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
}

type Config struct {
	Cors CorsConfig `json:"cors"`
}

var config *Config

func LoadConfig() error {
	// Get the directory of the current file
	_, filename, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(filename)

	// Read config file from the same directory as this file
	file, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	if err != nil {
		return err
	}

	config = &Config{}
	return json.Unmarshal(file, config)
}

func GetCorsConfig() *CorsConfig {
	return &config.Cors
}
