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

func (s *AssetsService) CreateAsset(asset *models.Assets) error {
	return mgm.Coll(asset).Create(asset)
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

func (s *AssetsService) UpdateAsset(id string, updateData *models.Assets) (*models.Assets, error) {
	asset, err := s.GetAssetByID(id)
	if err != nil {
		return nil, err
	}

	asset.Name = updateData.Name
	asset.URL = updateData.URL
	asset.UpdatedAt = time.Now()

	err = s.collection.Update(asset)
	return asset, err
}

func (s *AssetsService) DeleteAsset(id string) error {
	asset, err := s.GetAssetByID(id)
	if err != nil {
		return err
	}

	asset.DeletedAt = &time.Time{}
	return s.collection.Update(asset)
}
