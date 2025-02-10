package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gofast/routes"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	router.SetupRoutes(r)
	return r
}

func setupTestRoutes(t *testing.T) {
	r := setupTestRouter()

	t.Run("GET /health", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"message": "OK"}`, resp.Body.String())
	})

	t.Run("GET /api/v1", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"message": "Welcome to GoFast API v1"}`, resp.Body.String())
	})

	t.Run("GET /api/v1/hello", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/hello", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"message": "Hello, World!"}`, resp.Body.String())
	})
}