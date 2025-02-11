package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gofast/routes"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connexion à la base de données
	// database.Connect()
	// Connexion au service Azure Blob Storage
	client := getServiceClientTokenCredential()
	// fmt.Println("✅ Fichier uploadé avec succès !", client)
	fmt.Println("✅ Client create with success !", client)
	// containerName := "mycontainer"
	// blobName := "example.txt"
	// filepath := "path/files/file.txt"

	// // Uploader le fichier dans le Blob Storage
	// createContainer(client, containerName)
	// uploadBlobFile(client, containerName, blobName, filepath)

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	if err := r.Run(":80"); err != nil {
		log.Fatal("server startup failed:", err)
	}
}

func getServiceClientTokenCredential() *azblob.Client {
	// Variables (peuvent être définies dans des variables d'environnement)
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT") // Remplace par le nom du compte de stockage
	accountKey := os.Getenv("AZURE_STORAGE_KEY")      // Remplace par la clé d'accès obtenue via Azure CLI

	// Créer les credentials
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic(fmt.Sprintf("Erreur d'authentification: %v", err))
	}

	// Construire l'URL du service
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// Créer le client du service Azure Blob Storage
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		panic(fmt.Sprintf("Erreur de création du client de service: %v", err))
	}
	return client
}

func createContainer(client *azblob.Client, containerName string) {
	// Create a container
	_, err := client.CreateContainer(context.TODO(), containerName, nil)
	if err != nil {
		log.Fatalf("Error creating container : %s", err)
	}
}

func uploadBlobFile(client *azblob.Client, containerName string, blobName string, filepath string) {
	// Open the file for reading
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Error open file : %s", err)
	}

	defer file.Close()

	// Upload the file to the specified container with the specified blob name
	_, err = client.UploadFile(context.TODO(), containerName, blobName, file, nil)
	if err != nil {
		log.Fatalf("Error upload file on azure : %s", err)
	}
}
