package services_test

import (
	"testing"

	"gofast/services"
	"gofast/tests/helpers"
	"gofast/utils/validator"

	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	helpers.SetupTestDB(t)
	userService := services.NewUserService()

	t.Run("Register", func(t *testing.T) {
		t.Run("should register new user successfully", func(t *testing.T) {
			req := &services.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			}

			err := userService.Register(req)
			assert.NoError(t, err)
		})

		t.Run("should fail with invalid email", func(t *testing.T) {
			req := &services.RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
			}

			err := userService.Register(req)
			assert.ErrorIs(t, err, validator.ErrInvalidEmailFormat)
		})

		t.Run("should fail with duplicate email", func(t *testing.T) {
			req := &services.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			}

			err := userService.Register(req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "user already exists")
		})
	})

	t.Run("Login", func(t *testing.T) {
		t.Run("should login successfully", func(t *testing.T) {
			req := &services.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			}

			response, err := userService.Login(req)
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, "Login successful", response.Message)
			assert.Equal(t, req.Email, response.Email)
			assert.Equal(t, "USER", response.Role)
			assert.NotEmpty(t, response.Token)
			assert.NotEmpty(t, response.UserID)
		})

		t.Run("should fail with invalid credentials", func(t *testing.T) {
			req := &services.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			}

			response, err := userService.Login(req)
			assert.Error(t, err)
			assert.Nil(t, response)
			assert.Contains(t, err.Error(), "invalid credentials")
		})

		t.Run("should fail with non-existent user", func(t *testing.T) {
			req := &services.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			}

			response, err := userService.Login(req)
			assert.Error(t, err)
			assert.Nil(t, response)
			assert.Contains(t, err.Error(), "invalid credentials")
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		// First create a user to test with
		registerReq := &services.RegisterRequest{
			Email:    "getbyid@example.com",
			Password: "password123",
		}
		err := userService.Register(registerReq)
		assert.NoError(t, err)

		// Login to get the user ID
		loginReq := &services.LoginRequest{
			Email:    "getbyid@example.com",
			Password: "password123",
		}
		loginRes, err := userService.Login(loginReq)
		assert.NoError(t, err)

		t.Run("should get user by ID", func(t *testing.T) {
			user, err := userService.GetByID(loginRes.UserID)
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, registerReq.Email, user.Email)
			assert.Equal(t, "USER", user.Role)
		})

		t.Run("should fail with invalid ID", func(t *testing.T) {
			user, err := userService.GetByID("invalid-id")
			assert.Error(t, err)
			assert.Nil(t, user)
		})
	})
}
