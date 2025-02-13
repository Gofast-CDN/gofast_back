package services

import (
	"errors"
	"time"

	"gofast/models"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		URL:      url,
		Path:     filePath,
		ParentID: &parentAsset.ID,
		Childs:   []models.Assets{},
	}

	if err := mgm.Coll(asset).Create(asset); err != nil {
		return nil, err
	}

	if err := mgm.Coll(asset).FindByID(asset.ID.Hex(), asset); err != nil {
		return nil, err
	}

	parentAsset.Childs = append(parentAsset.Childs, *asset)
	parentAsset.Size += fileSize
	if err := mgm.Coll(parentAsset).Update(parentAsset); err != nil {
		return nil, errors.New("Impossible de mettre à jour le parent")
	}

	return asset, nil
}

func (s *AssetsService) CreateRepoAsset(id, containerName, blobName string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID invalide")
	}

	parentAsset, err := s.GetAssetByName(containerName)
	if err != nil {
		return errors.New("Impossible de retrouver le parent")
	}

	filePath := parentAsset.Path + "/" + blobName

	asset := &models.Assets{
		Name:     blobName,
		Type:     "folder",
		OwnerID:  userID,
		Size:     0,
		URL:      "",
		Path:     filePath,
		ParentID: &parentAsset.ID,
		Childs:   []models.Assets{},
	}

	if err := mgm.Coll(asset).Create(asset); err != nil {
		return err
	}

	if err := mgm.Coll(asset).FindByID(asset.ID.Hex(), asset); err != nil {
		return err
	}

	parentAsset.Childs = append(parentAsset.Childs, *asset)
	if err := mgm.Coll(parentAsset).Update(parentAsset); err != nil {
		return errors.New("Impossible de mettre à jour le parent")
	}

	return nil
}

func (s *AssetsService) CreateAsset(asset *models.Assets) error {
	return mgm.Coll(asset).Create(asset)
}

func (s *AssetsService) CreateRootRepoAsset(id string, repoName, repoPath string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID invalide")
	}

	// Create the asset object with empty parentId and childs
	asset := &models.Assets{
		Name:     repoName,
		Type:     "folder",
		OwnerID:  userID,
		Size:     0,
		URL:      "",
		Path:     repoPath,
		ParentID: nil,               // ParentID is nil
		Childs:   []models.Assets{}, // Childs is an empty slice
	}

	// Save the asset to the collection
	createErr := s.collection.Create(asset)
	if createErr != nil {
		return createErr
	}

	return nil
}

func (s *AssetsService) GetAllAssets() ([]models.Assets, error) {
	var assets []models.Assets
	err := s.collection.SimpleFind(&assets, bson.M{"deletedAt": nil})
	return assets, err
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
