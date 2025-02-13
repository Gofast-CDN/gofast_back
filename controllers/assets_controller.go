package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"gofast/models"
	"gofast/services"

	"github.com/gin-gonic/gin"
)

type AssetsController struct {
	assetsService *services.AssetsService
}

func NewAssetsController() *AssetsController {
	return &AssetsController{
		assetsService: services.NewAssetsService(),
	}
}

func (ctrl *AssetsController) CreateAsset(c *gin.Context) {
	userValue, _ := c.Get("user")
	user := userValue.(*models.User)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error while parsing the form for the file": err})
		return
	}
	defer file.Close() // Ensure we close the file after using it

	// Get original filename with extension
	originalFilename := fileHeader.Filename
	// Get file extension
	fileExt := filepath.Ext(originalFilename)

	containerName := c.DefaultPostForm("containerName", "default-container")
	blobName := c.DefaultPostForm("blobName", "default-blob-name")
	if !strings.HasSuffix(blobName, fileExt) {
		blobName = blobName + fileExt
	}
	fileSize := fileHeader.Size
	fmt.Println("Container: ", containerName, "; Blob:", blobName, ", Size: ", fileSize)

	// get repoassets by name
	repoAsset, err := ctrl.assetsService.GetAssetByName(containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les assets", "details": err.Error()})
		return
	}

	// connect to blob storage
	blobService, err := services.NewBlobStorageService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// upload file to blob storage into the good container
	fileURL, err := blobService.UploadFile(repoAsset.Name, blobName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// save file asset in db
	createdAssets, err := ctrl.assetsService.CreateFileAsset(containerName, blobName, fileURL, fileSize, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'asset"})
		return
	}

	// update parent repoassets in db
	repoAsset.Childs = append(repoAsset.Childs, createdAssets.ID)
	_, err = ctrl.assetsService.UpdateAsset(repoAsset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de mettre à jour l'asset parent"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset créé avec succès"})
}

func (ctrl *AssetsController) CreateFolder(c *gin.Context) {
	var asset models.Assets

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.assetsService.CreateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer le dossier"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Dossier créé avec succès", "data": asset})
}

func (ctrl *AssetsController) GetAssets(c *gin.Context) {
	userValue, _ := c.Get("user")
	user := userValue.(*models.User)
	createRootRepoAssetName := user.ID.Hex() + "-root"

	fmt.Println("Assets root:", createRootRepoAssetName)
	fmt.Println("Assets root:", createRootRepoAssetName)

	assets, err := ctrl.assetsService.GetAssetByName(createRootRepoAssetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les assets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assets)
}

func (ctrl *AssetsController) GetAssetByID(c *gin.Context) {
	id := c.Param("id")

	asset, err := ctrl.assetsService.GetAssetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset non trouvé"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (ctrl *AssetsController) UpdateAsset(c *gin.Context) {
	id := c.Param("id")

	var updateData models.Assets
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get assets by id
	assetToUpdate, err := ctrl.assetsService.GetAssetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer l'asset", "details": err.Error()})
		return
	}

	// update asset
	asset, err := ctrl.assetsService.UpdateAsset(assetToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset mis à jour avec succès", "data": asset})
}

func (ctrl *AssetsController) DeleteAsset(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.assetsService.DeleteAsset(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset supprimé avec succès"})
}
