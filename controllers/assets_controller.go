package controllers

import (
	"fmt"
	"net/http"

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
	var asset models.Assets

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Asset : ", asset)

	file, _, err := c.Request.FormFile("file") // Assuming the form field is "file"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	fmt.Println("File : ", file)

	// You may need to extract containerName and blobName from the form data as well
	containerName := c.DefaultPostForm("containerName", "default-container")
	blobName := c.DefaultPostForm("blobName", "default-blob-name")

	blobService, err := services.NewBlobStorageService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := blobService.UploadFile(containerName, blobName, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.assetsService.CreateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'asset"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset créé avec succès", "data": asset})
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
	assets, err := ctrl.assetsService.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les assets"})
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

	asset, err := ctrl.assetsService.UpdateAsset(id, &updateData)
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
