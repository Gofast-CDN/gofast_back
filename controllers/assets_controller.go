package controllers

import (
	"net/http"

	"gofast/models"
	"gofast/services"

	"github.com/gin-gonic/gin"
)

type AssetsController struct {
	assetsService *services.AssetsService
}

func NewAssetsController(blobService services.BlobStorage) *AssetsController {
	return &AssetsController{
		assetsService: services.NewAssetsService(blobService),
	}
}

// 🔹 Création d’un asset
func (ctrl *AssetsController) CreateAsset(c *gin.Context) {
	var asset models.Assets
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ctrl.assetsService.CreateAsset(&asset, "./test/test-file.pdf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'asset"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Asset créé avec succès", "data": asset})
}

// 🔹 Récupération d’un asset par ID
func (ctrl *AssetsController) GetAssetByID(c *gin.Context) {
	id := c.Param("id")
	asset, err := ctrl.assetsService.GetAssetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset non trouvé"})
		return
	}
	c.JSON(http.StatusOK, asset)
}

// 🔹 Suppression d’un asset
func (ctrl *AssetsController) DeleteAsset(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.assetsService.DeleteAsset(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset supprimé avec succès"})
}
