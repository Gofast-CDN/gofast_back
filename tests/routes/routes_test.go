package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"gofast/routes"
	"gofast/services"
	"gofast/tests/helpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type PartialHealthResponse struct {
	Status    string `json:"status"`
	GoVersion string `json:"go_version"`
}

func setupTestRouter(t *testing.T) *gin.Engine {
	// Setup test database
	helpers.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

func TestSetupRoutes(t *testing.T) {
	r := setupTestRouter(t)

	t.Run("Health Check Routes", func(t *testing.T) {
		t.Run("GET /health", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response PartialHealthResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "healthy", response.Status)
			assert.Equal(t, runtime.Version(), response.GoVersion)
		})
	})
	t.Run("User Routes", func(t *testing.T) {
		t.Run("POST /api/v1/users - Create User", func(t *testing.T) {
			// Create user test
			userData := services.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			}
			jsonData, err := json.Marshal(userData)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, userData.Email, response["email"])
			assert.Equal(t, "USER", response["role"])
			assert.NotEmpty(t, response["token"])
		})

		t.Run("POST /api/v1/users - Duplicate Email", func(t *testing.T) {
			// Test duplicate email
			userData := services.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			}
			jsonData, _ := json.Marshal(userData)

			req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})

}
