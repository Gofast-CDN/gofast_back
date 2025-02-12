package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type BlobStorageService struct {
	Client *azblob.Client
}

func NewBlobStorageService() (*BlobStorageService, error) {
	client, err := getServiceClientTokenCredential()
	if err != nil {
		return nil, fmt.Errorf("Failed to create blob storage client: %v", err)
	}

	return &BlobStorageService{
		Client: client,
	}, nil
}

func getServiceClientTokenCredential() (*azblob.Client, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	// Vérification des variables d'environnement
	if accountName == "" || accountKey == "" {
		return nil, fmt.Errorf("❌ Environnement variables AZURE_STORAGE_ACCOUNT_NAME or AZURE_STORAGE_ACCOUNT_KEY are not provide")
	}

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create credentials: %v", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	fmt.Println("✅ Blob Storage service initialized successfully!")
	return client, nil
}

// Ajouter des méthodes pour gérer les opérations de blob storage
func (service *BlobStorageService) UploadFile(containerName, blobName string, file multipart.File) error {

	_, err := service.Client.UploadStream(
		context.TODO(),
		containerName,
		blobName,
		file,
		nil,
	)

	if err != nil {
		log.Fatalf("Error uploading file on azure: %s", err)
	}
	return nil
}

func (service *BlobStorageService) CreateContainer(containerName string) (azblob.CreateContainerResponse, error) {
	// Create a container
	response, err := service.Client.CreateContainer(context.TODO(), containerName, nil)
	if err != nil {
		log.Fatalf("Error creating container on azure: %s", err)
	}

	return response, nil
}

func (service *BlobStorageService) DeleteContainer(containerName string) error {
	// Delete the container
	_, err := service.Client.DeleteContainer(context.TODO(), containerName, nil)
	if err != nil {
		log.Fatalf("Error deleting container on Azure: %s", err)
		return err // Return the error for further handling if necessary
	}

	return nil
}
