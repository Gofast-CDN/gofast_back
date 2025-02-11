package services

import (
	"net/http"
	"time"

	"gofast/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Créer un fichier ou un dossier
func CreateAsset(c *gin.Context) {
	var asset models.Assets

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mgm.Coll(&asset).Create(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'asset"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset créé avec succès", "data": asset})
}

// Récupérer tous les fichiers et dossiers
func GetAssets(c *gin.Context) {
	var assets []models.Assets

	// Récupérer tous les fichiers/dossiers non supprimés (soft delete)
	err := mgm.Coll(&models.Assets{}).SimpleFind(&assets, bson.M{"deletedAt": nil})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les assets"})
		return
	}

	c.JSON(http.StatusOK, assets)
}

// Récupérer un fichier/dossier par ID
func GetAssetByID(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var asset models.Assets

	err = mgm.Coll(&models.Assets{}).First(bson.M{"_id": objectID, "deletedAt": nil}, &asset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset non trouvé"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

// Mettre à jour un fichier/dossier
func UpdateAsset(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var updateData models.Assets
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var asset models.Assets
	err = mgm.Coll(&models.Assets{}).First(bson.M{"_id": objectID, "deletedAt": nil}, &asset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset non trouvé"})
		return
	}

	asset.Name = updateData.Name
	asset.URL = updateData.URL
	asset.UpdatedAt = time.Now()

	if err := mgm.Coll(&asset).Update(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset mis à jour avec succès", "data": asset})
}

// Supprimer un fichier/dossier (Soft Delete)
func DeleteAsset(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var asset models.Assets
	err = mgm.Coll(&models.Assets{}).First(bson.M{"_id": objectID, "deletedAt": nil}, &asset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset non trouvé"})
		return
	}

	// soft delete
	asset.DeletedAt = time.Now().Format(time.RFC3339)

	if err := mgm.Coll(&asset).Update(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset supprimé avec succès"})
}
