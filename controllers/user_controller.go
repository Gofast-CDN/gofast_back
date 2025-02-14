package controllers

import (
	"fmt"
	"net/http"

	"gofast/models"
	"gofast/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

func (uc *UserController) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse, err := uc.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID := userResponse.UserID
	rootRepoName := userID + "-root"

	blobService, err := services.NewBlobStorageService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := blobService.CreateContainer(rootRepoName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	assetsService := services.NewAssetsService()

	rootRepoPath := "/" + rootRepoName

	rootContainerID, err := assetsService.CreateRootRepoAsset(userID, rootRepoName, rootRepoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update user with root container ID
	_, err = uc.userService.UpdateUserRootRepoAsset(userID, rootContainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})

}

func (uc *UserController) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (uc *UserController) GetMe(c *gin.Context) {
	userValue, _ := c.Get("user")
	user := userValue.(*models.User)
	c.JSON(http.StatusOK, gin.H{
		"id":              user.ID,
		"email":           user.Email,
		"role":            user.Role,
		"rootContainerID": user.RootContainerID,
	})
}

func (uc *UserController) Delete(c *gin.Context) {
	// Extract userID from request URL or JSON body
	userID := c.Param("userId") // If userId is passed as a URL parameter

	// Delete user by ID
	if err := uc.userService.DeleteUserByID(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rootContainerName := userID + "-root"
	fmt.Println("User ID retrieved from context:", userID)

	blobService, err := services.NewBlobStorageService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := blobService.DeleteContainer(rootContainerName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
