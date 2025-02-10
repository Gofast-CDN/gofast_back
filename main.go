package main

import (
    "gofast/routes"

    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize Gin router
    r := gin.Default()

    // Setup routes
    routes.SetupRoutes(r)

    // Start server
    r.Run(":8080")
}
