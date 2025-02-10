package handlers

import (
    "net/http"
    "runtime"
    "time"

    "github.com/gin-gonic/gin"
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
