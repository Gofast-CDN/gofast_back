package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gofast/middleware"
	"gofast/utils/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		email, _ := c.Get("email")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{
			"userID": userID,
			"email":  email,
			"role":   role,
		})
	})
	return r
}

func TestAuthMiddleware(t *testing.T) {
	r := setupTestRouter()

	t.Run("should allow access with valid token", func(t *testing.T) {
		token, _ := auth.GenerateToken("123", "test@example.com", "USER")

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "123", response["userID"])
		assert.Equal(t, "test@example.com", response["email"])
		assert.Equal(t, "USER", response["role"])
	})

	t.Run("should fail without Authorization header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should fail with invalid token format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
