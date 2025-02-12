package services

import (
	"context"
	"fmt"
	"log"
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
func (service *BlobStorageService) UploadFile(containerName, blobName, filepath string) (azblob.UploadFileResponse, error) {
	// Open the file for reading
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return azblob.UploadFileResponse{}, fmt.Errorf("Error open file: %v", err)
	}

	defer file.Close()

	// Upload the file to the specified container with the specified blob name
	response, err := service.Client.UploadFile(context.TODO(), containerName, blobName, file, nil)
	if err != nil {
		return azblob.UploadFileResponse{}, fmt.Errorf("Error upload file on azure: %s", err)
	}

	// Implémentation de l'upload
	return response, nil
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
