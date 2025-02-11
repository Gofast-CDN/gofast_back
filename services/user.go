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

// Créer un utilisateur
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'utilisateur existe déjà
	var existingUser models.User
	err := mgm.Coll(&models.User{}).First(bson.M{"email": user.Email}, &existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Cet utilisateur existe déjà"})
		return
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Créer l'utilisateur
	if err := mgm.Coll(&user).Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Utilisateur créé avec succès", "user": user})
}

// Récupérer tous les utilisateurs (sans les supprimés)
func GetUsers(c *gin.Context) {
	var users []models.User

	// Récupérer tous les utilisateurs non supprimés (soft delete)
	err := mgm.Coll(&models.User{}).SimpleFind(&users, bson.M{"deletedAt": nil})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les utilisateurs"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Récupérer un utilisateur par ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	// Convertir en ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var user models.User
	err = mgm.Coll(&models.User{}).First(bson.M{"_id": objectID, "deletedAt": nil}, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Mettre à jour un utilisateur
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Convertir en ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trouver l'utilisateur par ID
	var user models.User
	err = mgm.Coll(&models.User{}).First(bson.M{"_id": objectID, "deletedAt": nil}, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Mise à jour des données si elles sont fournies
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Role != "" {
		user.Role = updateData.Role
	}

	user.UpdatedAt = time.Now()

	// Sauvegarder les modifications
	if err := mgm.Coll(&user).Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur mis à jour avec succès", "user": user})
}

// Supprimer un utilisateur (Soft Delete)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Convertir en ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var user models.User
	err = mgm.Coll(&models.User{}).First(bson.M{"_id": objectID, "deletedAt": nil}, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Soft delete (marquer comme supprimé)
	now := time.Now()
	user.DeletedAt = &now

	if err := mgm.Coll(&user).Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur supprimé avec succès"})
}
