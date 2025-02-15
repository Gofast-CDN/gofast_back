package services

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
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
		return nil, fmt.Errorf("Environnement variables AZURE_STORAGE_ACCOUNT_NAME or AZURE_STORAGE_ACCOUNT_KEY are not provide")
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
func (service *BlobStorageService) UploadFile(containerName, blobName string, file multipart.File) (string, error) {
	_, err := service.Client.UploadStream(
		context.TODO(),
		containerName,
		blobName,
		file,
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("Error uploading file: %w", err) // Ne pas utiliser log.Fatalf
	}

	// Construct the Blob URL
	blobURL, err := service.GetBlobSASURL(containerName, blobName)
	if err != nil {
		return "", fmt.Errorf("failed to get the SAS URL: %w", err)
	}

	return blobURL, nil
}

func (service *BlobStorageService) GetContentTypeFromFile(blobName string) string {
	switch {
	case strings.HasSuffix(blobName, ".jpg"), strings.HasSuffix(blobName, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(blobName, ".png"):
		return "image/png"
	case strings.HasSuffix(blobName, ".gif"):
		return "image/gif"
	case strings.HasSuffix(blobName, ".webp"):
		return "image/webp"
	case strings.HasSuffix(blobName, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(blobName, ".pdf"):
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func (service *BlobStorageService) GetBlobSASURL(containerName, blobName string) (string, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") // Get Storage Account Name
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")   // Get Storage Account Key

	// Validate credentials
	if accountName == "" || accountKey == "" {
		return "", fmt.Errorf("AZURE_STORAGE_ACCOUNT_NAME or AZURE_STORAGE_ACCOUNT_KEY not set")
	}

	// Generate SharedKeyCredential for SAS token signing
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", fmt.Errorf("failed to create shared key credential: %w", err)
	}

	// Set SAS token parameters
	startTime := time.Now().UTC()
	expiryTime := startTime.Add(24 * time.Hour) // 24-hour validity

	permissions := sas.BlobPermissions{
		Read: true,
		List: true, // Ajouter la permission de lister
	}
	sasValues := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS, // Restrict to HTTPS
		StartTime:     startTime,         // Start time (optional)
		ExpiryTime:    expiryTime,        // Expiration time
		Permissions:   permissions.String(),
		ContainerName: containerName,
		BlobName:      blobName,
		ContentType:   "inline",                   // Ajouter cette ligne
		CacheControl:  "public, max-age=31536000", // Cache d'un an
	}

	// Generate SAS query parameters
	sasQueryParams, err := sasValues.SignWithSharedKey(cred)
	if err != nil {
		return "", fmt.Errorf("failed to sign SAS token: %w", err)
	}

	// Construct final URL
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		accountName,
		containerName,
		blobName,
		sasQueryParams.Encode(),
	)

	return blobURL, nil
}

// Ajouter des méthodes pour gérer les opérations de blob storage
func (service *BlobStorageService) DownloadFile(containerName, blobName string) (bytes.Buffer, error) {
	// Download the blob
	get, err := service.Client.DownloadStream(context.TODO(), containerName, blobName, nil)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("Error downloading blob from azure: %s", err)
	}

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(context.TODO(), &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("Error reading download blob: %s", err)
	}

	err = retryReader.Close()
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("Error closing reader: %s", err)
	}

	// Print the contents of the blob we created
	fmt.Println("Blob contents:")
	fmt.Println(downloadedData.String())

	return downloadedData, nil
}

func (service *BlobStorageService) CreateContainer(containerName string) (azblob.CreateContainerResponse, error) {
	// Create a container
	response, err := service.Client.CreateContainer(context.TODO(), containerName, nil)
	if err != nil {
		return azblob.CreateContainerResponse{}, fmt.Errorf("Error creating container on azure: %s", err)
	}

	return response, nil
}

func (service *BlobStorageService) DeleteBlob(containerName, blobName string) (azblob.DeleteBlobResponse, error) {
	// Create a container
	response, err := service.Client.DeleteBlob(context.TODO(), containerName, blobName, nil)
	if err != nil {
		return azblob.DeleteBlobResponse{}, fmt.Errorf("Error deleting blob on azure: %s", err)
	}

	return response, nil
}

func (service *BlobStorageService) DeleteContainer(containerName string) (azblob.DeleteContainerResponse, error) {
	// Create a container
	response, err := service.Client.DeleteContainer(context.TODO(), containerName, nil)
	if err != nil {
		return azblob.DeleteContainerResponse{}, fmt.Errorf("Error deleting container on azure: %s", err)
	}

	return response, nil
}

func (service *BlobStorageService) GetContainerByContainerName(containerName string) (*runtime.Pager[azblob.ListContainersResponse], error) {
	// List the containers in the storage account with a prefix
	pager := service.Client.NewListContainersPager(&azblob.ListContainersOptions{
		Prefix: &containerName,
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("Error deleting container on azure: %s", err)
		}

		for _, container := range resp.ContainerItems {
			fmt.Println(*container.Name)
		}
	}

	return pager, nil
}
