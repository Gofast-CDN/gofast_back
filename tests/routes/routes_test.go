package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gofast/config"
	"gofast/routes"
	"gofast/services"
	"gofast/tests/helpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	if err := config.LoadConfig(); err != nil {
		t.Fatal("Failed to load config:", err)
	}

	helpers.SetupTestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

func TestSetupRoutes(t *testing.T) {
	r := setupTestRouter(t)

	t.Run("Auth Routes", func(t *testing.T) {
		t.Run("POST /api/v1/users/register - Register User", func(t *testing.T) {
			userData := services.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			}
			jsonData, err := json.Marshal(userData)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Registration successful", response["message"])
		})

		t.Run("POST /api/v1/users/login - Login Success", func(t *testing.T) {
			loginData := services.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			}
			jsonData, _ := json.Marshal(loginData)

			req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Login successful", response["message"])
			assert.Equal(t, loginData.Email, response["email"])
			assert.NotEmpty(t, response["token"])
		})

		t.Run("GET /api/v1/users/me - Get Profile (Protected)", func(t *testing.T) {
			// First login to get token
			loginData := services.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			}
			jsonData, _ := json.Marshal(loginData)

			loginReq := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(jsonData))
			loginReq.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, loginReq)

			var loginResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
			assert.NoError(t, err, "Failed to unmarshal login response")
			token := loginResponse["token"].(string)

			// Test protected route
			req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w = httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, loginData.Email, response["email"])
			assert.Equal(t, "USER", response["role"])
		})
		t.Run("GET /api/v1/users/me - Unauthorized", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("POST /api/v1/users/login - Invalid Credentials", func(t *testing.T) {
			loginData := services.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			}
			jsonData, _ := json.Marshal(loginData)

			req := httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	})
}
