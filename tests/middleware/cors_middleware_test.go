package middleware_test

import (
	"net/http/httptest"
	"testing"

	"gofast/config"
	"gofast/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	// Load config before tests
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CorsMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})
	return r
}

func TestCorsMiddleware(t *testing.T) {
	router := setupRouter()

	t.Run("should allow configured origin", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")

		router.ServeHTTP(w, req)

		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, 200, w.Code)
	})

	t.Run("should not allow unconfigured origin", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://malicious.com")

		router.ServeHTTP(w, req)

		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, 200, w.Code)
	})

	t.Run("should handle OPTIONS request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")

		router.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("should set correct CORS headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")

		router.ServeHTTP(w, req)

		corsConfig := config.GetCorsConfig()
		expectedMethods := "GET, POST, PUT, DELETE, OPTIONS"
		expectedHeaders := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"

		assert.Equal(t, expectedMethods, w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, expectedHeaders, w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, corsConfig.AllowCredentials, w.Header().Get("Access-Control-Allow-Credentials") == "true")
	})
}
