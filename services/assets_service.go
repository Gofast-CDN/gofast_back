package services

import (
	"errors"
	"fmt"
	"time"

	"gofast/models"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssetsService struct {
	collection *mgm.Collection
}

func NewAssetsService() *AssetsService {
	return &AssetsService{
		collection: mgm.Coll(&models.Assets{}),
	}
}

func (s *AssetsService) CreateFileAsset(containerName, blobName, url string, fileSize int64, userID primitive.ObjectID) (*models.Assets, error) {
	parentAsset, err := s.GetAssetByName(containerName)
	if err != nil {
		return nil, errors.New("Impossible de retrouver le parent")
	}

	filePath := parentAsset.Path + "/" + blobName

	asset := &models.Assets{
		Name:     blobName,
		Type:     "file",
		OwnerID:  userID,
		Size:     fileSize,
		Depth:    0,
		URL:      url,
		Path:     filePath,
		ParentID: &parentAsset.ID,
		Childs:   []models.Assets{},
	}

	if err := mgm.Coll(asset).Create(asset); err != nil {
		return nil, err
	}

	// Updating the path after creating the asset to get its id
	filePathInfo := models.PathInfoEntry{
		ContainerID:   asset.ID.Hex(),
		ContainerName: asset.Name,
	}

	newPathInfo := append(parentAsset.PathInfo, filePathInfo)
	asset.PathInfo = newPathInfo
	if err := mgm.Coll(asset).Update(asset); err != nil {
		return nil, errors.New("Impossible de mettre à jour l'asset")
	}

	parentAsset.Childs = append(parentAsset.Childs, *asset)
	parentAsset.NbChildren++
	parentAsset.Size += fileSize
	if err := mgm.Coll(parentAsset).Update(parentAsset); err != nil {
		return nil, errors.New("Impossible de mettre à jour le parent")
	}

	return asset, nil
}

func (s *AssetsService) CreateRepoAsset(userID primitive.ObjectID, newContainerName, parentID string) error {
	parentAsset, err := s.GetAssetByID(parentID)
	if err != nil {
		return errors.New("Impossible de retrouver le parent")
	}

	if parentAsset.Depth >= 10 {
		return errors.New("Profondeur maximale atteinte")
	}

	filePath := parentAsset.Path + "/" + newContainerName

	asset := &models.Assets{
		Name:     newContainerName,
		Type:     "folder",
		OwnerID:  userID,
		Size:     0,
		Depth:    parentAsset.Depth + 1,
		URL:      "",
		Path:     filePath,
		ParentID: &parentAsset.ID,
		Childs:   []models.Assets{},
	}

	if err := mgm.Coll(asset).Create(asset); err != nil {
		return err
	}

	// Updating the path after creating the asset to get its id
	filePathInfo := models.PathInfoEntry{
		ContainerID:   asset.ID.Hex(),
		ContainerName: asset.Name,
	}
	newPathInfo := append(parentAsset.PathInfo, filePathInfo)
	asset.PathInfo = newPathInfo
	if err := mgm.Coll(asset).Update(asset); err != nil {
		return errors.New("Impossible de mettre à jour l'asset")
	}

	parentAsset.Childs = append(parentAsset.Childs, *asset)
	parentAsset.NbChildren++
	if err := mgm.Coll(parentAsset).Update(parentAsset); err != nil {
		return errors.New("Impossible de mettre à jour le parent")
	}

	return nil
}

func (s *AssetsService) CreateAsset(asset *models.Assets) error {
	return mgm.Coll(asset).Create(asset)
}

func (s *AssetsService) CreateRootRepoAsset(id string, repoName, repoPath string) (primitive.ObjectID, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, errors.New("ID invalide")
	}

	// Create the asset object with empty parentId and childs
	asset := &models.Assets{
		Name:     repoName,
		Type:     "folder",
		OwnerID:  userID,
		Size:     0,
		Depth:    0,
		URL:      "",
		Path:     repoPath,
		ParentID: nil,               // ParentID is nil
		Childs:   []models.Assets{}, // Childs is an empty slice
	}

	// Save the asset to the collection
	createErr := s.collection.Create(asset)
	if createErr != nil {
		return primitive.NilObjectID, createErr
	}

	// Updating the path after creating the asset to get its id
	filePathInfo := models.PathInfoEntry{
		ContainerID:   asset.ID.Hex(),
		ContainerName: "home",
	}
	asset.PathInfo = []models.PathInfoEntry{filePathInfo}

	if err := mgm.Coll(asset).Update(asset); err != nil {
		return primitive.NilObjectID, errors.New("Impossible de mettre à jour l'asset")
	}

	return asset.ID, nil
}

func (s *AssetsService) GetRecentUserAssetsFiles(userID primitive.ObjectID) ([]models.Assets, error) {
	var assets []models.Assets

	err := s.collection.SimpleFind(&assets, bson.M{
		"ownerId":   userID,
		"type":      "file",
		"deletedAt": nil,
	}, &options.FindOptions{
		Sort:  bson.M{"updated_at": -1},
		Limit: func(i int64) *int64 { return &i }(10),
	})

	if err != nil {
		return nil, err
	}

	if assets == nil {
		assets = make([]models.Assets, 0)
	}

	return assets, nil
}

func (s *AssetsService) GetRecentUserAssetsFolder(userID primitive.ObjectID) ([]models.Assets, error) {
	var assets []models.Assets

	err := s.collection.SimpleFind(&assets, bson.M{
		"ownerId":   userID,
		"type":      "folder",
		"deletedAt": nil,
	}, &options.FindOptions{
		Sort:  bson.M{"updated_at": -1},
		Limit: func(i int64) *int64 { return &i }(10),
	})

	if err != nil {
		return nil, err
	}

	if assets == nil {
		assets = make([]models.Assets, 0)
	}

	return assets, nil
}

func (s *AssetsService) GetUserAssetByID(id string, userID primitive.ObjectID) (*models.Assets, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID invalide")
	}

	fmt.Println("id:", objectID, "userID:", userID)

	var asset models.Assets
	err = s.collection.First(bson.M{"_id": objectID, "ownerId": userID, "deletedAt": nil}, &asset)
	return &asset, err
}

func (s *AssetsService) GetAssetByID(id string) (*models.Assets, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID invalide")
	}

	var asset models.Assets
	err = s.collection.First(bson.M{"_id": objectID, "deletedAt": nil}, &asset)
	return &asset, err
}

func (s *AssetsService) GetAssetByName(name string) (*models.Assets, error) {
	var asset models.Assets
	err := s.collection.First(bson.M{"name": name, "deletedAt": nil}, &asset)
	if err != nil {
		return nil, errors.New("Asset not found")
	}
	return &asset, nil
}

func (s *AssetsService) UpdateAsset(updateAsset *models.Assets) (*models.Assets, error) {

	updateAsset.UpdatedAt = time.Now()

	err := s.collection.Update(updateAsset)
	if err != nil {
		return nil, err
	}

	return updateAsset, nil
}

func (s *AssetsService) DeleteAsset(id string) error {
	asset, err := s.GetAssetByID(id)
	if err != nil {
		return err
	}

	asset.DeletedAt = &time.Time{}
	return s.collection.Update(asset)
}
