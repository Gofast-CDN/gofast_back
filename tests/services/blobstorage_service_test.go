package services_test

import (
	"os"
	"testing"

	"gofast/services"

	"github.com/stretchr/testify/assert"
)

func TestBlobstorageService(t *testing.T) {
	// Setup - Définir les variables d'environnement pour les tests
	os.Setenv("AZURE_STORAGE_ACCOUNT_NAME", "testaccount")
	os.Setenv("AZURE_STORAGE_ACCOUNT_KEY", "testkey")

	t.Run("should fail with invalid credentials", func(t *testing.T) {
		os.Setenv("AZURE_STORAGE_ACCOUNT_NAME", "")
		os.Setenv("AZURE_STORAGE_ACCOUNT_KEY", "")

		_, err := services.NewBlobStorageService()
		assert.Error(t, err)
	})
}
