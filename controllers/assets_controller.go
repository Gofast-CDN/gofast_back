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

func (ctrl *AssetsController) CreateFileAsset(c *gin.Context) {
	userValue, _ := c.Get("user")
	user := userValue.(*models.User)
	rootRepoAssetName := user.ID.Hex() + "-root"

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

	containerID := c.DefaultPostForm("containerId", "default-id")
	blobName := c.DefaultPostForm("blobName", "default-blob-name")
	if !strings.HasSuffix(blobName, fileExt) {
		blobName = blobName + fileExt
	}
	fileSize := fileHeader.Size

	// get repoassets by id and owner
	repoAsset, err := ctrl.assetsService.GetUserAssetByID(containerID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Immpossible de récupérer le dossier", "details": err.Error()})
		return
	}

	// connect to blob storage
	blobService, err := services.NewBlobStorageService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// upload file to blob storage into the good container
	fileURL, err := blobService.UploadFile(rootRepoAssetName, blobName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// save file asset in db
	_, err = ctrl.assetsService.CreateFileAsset(repoAsset.Name, blobName, fileURL, fileSize, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'asset"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset créé avec succès"})
}

type CreateRepoAssetRequest struct {
	ParentID      string `json:"containerId"`
	ContainerName string `json:"containerName"`
}

func (ctrl *AssetsController) CreateRepoAsset(c *gin.Context) {
	userValue, _ := c.Get("user")
	user := userValue.(*models.User)

	var req CreateRepoAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.assetsService.CreateRepoAsset(user.ID, req.ContainerName, req.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'asset"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Dossier créé avec succès"})
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
