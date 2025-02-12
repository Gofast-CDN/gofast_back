package services

import (
	"errors"
	"time"

	"gofast/models"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlobStorage interface {
	UploadFile(containerName, blobName, filePath string) error
	DeleteBlob(containerName, blobName string) error
}

type AssetsService struct {
	blobService BlobStorage
	collection  *mgm.Collection
}

func NewAssetsService(blobService BlobStorage) *AssetsService {
	return &AssetsService{
		blobService: blobService,
		collection:  mgm.Coll(&models.Assets{}),
	}
}

func (s *AssetsService) CreateAsset(asset *models.Assets, filePath string) error {
	err := s.blobService.UploadFile("assets", asset.ID.Hex()+"-"+asset.Name, filePath)
	if err != nil {
		return err
	}

	asset.URL = "https://storage.blob.core.windows.net/assets/" + asset.ID.Hex() + "-" + asset.Name
	return s.collection.Create(asset)
}

func (s *AssetsService) GetAssetByID(id string) (*models.Assets, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID invalide")
	}

	var asset models.Assets
	err = s.collection.First(bson.M{"_id": objectID, "deletedAt": nil}, &asset)
	if err != nil {
		return nil, errors.New("Asset non trouvé")
	}

	return &asset, nil
}

func (s *AssetsService) DeleteAsset(id string) error {
	asset, err := s.GetAssetByID(id)
	if err != nil {
		return err
	}

	err = s.blobService.DeleteBlob("assets", asset.ID.Hex()+"-"+asset.Name)
	if err != nil {
		return err
	}

	now := time.Now()
	asset.DeletedAt = &now
	return s.collection.Update(asset)
}
