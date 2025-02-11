package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"gofast/database"
)

var startTime = time.Now()

type HealthResponse struct {
	Status    string        `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Uptime    time.Duration `json:"uptime"`
	GoVersion string        `json:"go_version"`
}

func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).Round(time.Second),
		GoVersion: runtime.Version(),
	}

	c.JSON(http.StatusOK, response)
}

func MongoDBHealthCheck(c *gin.Context) {
	// Tenter de pinger la base de données
	err := database.Client.Ping(c, nil)
	if err != nil {
		// Si une erreur survient, la connexion n'est pas saine
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "message": "MongoDB connection failed"})
		return
	}

	// Si aucune erreur, la connexion est saine
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "message": "MongoDB connection is healthy"})
}
