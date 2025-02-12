package services_test

import (
	"testing"

	"gofast/services"
	"gofast/tests/helpers"

	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	helpers.SetupTestDB(t)
	userService := services.NewUserService()

	t.Run("Create User", func(t *testing.T) {
		t.Run("should create user with hashed password", func(t *testing.T) {
			req := &services.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			}

			response, err := userService.Create(req)
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, req.Email, response.Email)
			assert.Equal(t, "USER", response.Role)
			assert.NotEmpty(t, response.Token)
		})

		t.Run("should fail with duplicate email", func(t *testing.T) {
			req := &services.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			}

			response, err := userService.Create(req)
			assert.Error(t, err)
			assert.Nil(t, response)
			assert.Contains(t, err.Error(), "user already exists")
		})

		t.Run("should fail with invalid email", func(t *testing.T) {
			req := &services.CreateUserRequest{
				Email:    "invalid-email",
				Password: "password123",
			}

			response, err := userService.Create(req)
			assert.Error(t, err)
			assert.Nil(t, response)
		})
	})
}
