package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gofast/routes"
)

type PartialHealthResponse struct {
	Status    string `json:"status"`
	GoVersion string `json:"go_version"`
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

func TestSetupRoutes(t *testing.T) {
	r := setupTestRouter()

	t.Run("GET /health", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var actual PartialHealthResponse
		err := json.Unmarshal(resp.Body.Bytes(), &actual)
		assert.NoError(t, err)

		assert.Equal(t, "healthy", actual.Status)
		assert.Equal(t, runtime.Version(), actual.GoVersion)
	})

	t.Run("GET /api/v1", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"message": "Welcome to GoFast API v1"}`, resp.Body.String())
	})

	/**
	 * Users routes testing
	 */
	t.Run("POST /api/v1/users", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/users", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error": "EOF"}`, resp.Body.String())
	})
	t.Run("GET /api/v1/users", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `[]`, resp.Body.String())
	})

}
