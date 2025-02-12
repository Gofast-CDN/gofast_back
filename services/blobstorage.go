package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type BlobStorageService struct {
	client *azblob.Client
}

func NewBlobStorageService() (*BlobStorageService, error) {
	client, err := getServiceClientTokenCredential()
	if err != nil {
		return nil, fmt.Errorf("Failed to create blob storage client: %v", err)
	}

	fmt.Println("✅ Blob Storage service initialized successfully!")

	return &BlobStorageService{
		client: client,
	}, nil
}

func getServiceClientTokenCredential() (*azblob.Client, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create credentials: %v", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	return client, nil
}

// Ajouter des méthodes pour gérer les opérations de blob storage
func (service *BlobStorageService) UploadFile(containerName, blobName, filepath string) (azblob.UploadFileResponse, error) {
	// Open the file for reading
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return azblob.UploadFileResponse{}, fmt.Errorf("Error open file: %v", err)
	}

	defer file.Close()

	// Upload the file to the specified container with the specified blob name
	response, err := service.client.UploadFile(context.TODO(), containerName, blobName, file, nil)
	if err != nil {
		return azblob.UploadFileResponse{}, fmt.Errorf("Error upload file on azure: %s", err)
	}

	// Implémentation de l'upload
	return response, nil
}

func (service *BlobStorageService) createContainer(containerName string) (azblob.CreateContainerResponse, error) {
	// Create a container
	response, err := service.client.CreateContainer(context.TODO(), containerName, nil)
	if err != nil {
		log.Fatalf("Error creating container on azure: %s", err)
	}

	return response, nil
}
