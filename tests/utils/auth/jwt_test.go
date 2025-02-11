package auth_test

import (
	"os"
	"testing"

	"gofast/utils/auth"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key")

	t.Run("should generate valid token", func(t *testing.T) {
		userID := "123"
		email := "test@example.com"
		role := "USER"

		token, err := auth.GenerateToken(userID, email, role)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := auth.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		_, err := auth.ValidateToken("invalid-token")
		assert.Error(t, err)
	})
}
