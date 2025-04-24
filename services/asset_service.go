package services

import (
	"asset-maintenance/models"
	"asset-maintenance/repositories"
	"errors"
)

type AssetService interface {
	GetAllAssets() ([]models.Asset, error)
	GetAssetByID(id string) (*models.Asset, error)
	CreateAsset(asset *models.Asset, userID uint) error
	UpdateAsset(id string, input models.Asset) (*models.Asset, error)
	DeleteAsset(id string, userID uint) error
}

type assetService struct {
	assetRepo repositories.AssetRepository
}

func NewAssetService(assetRepo repositories.AssetRepository) AssetService {
	return &assetService{assetRepo: assetRepo}
}

func (s *assetService) GetAllAssets() ([]models.Asset, error) {
	return s.assetRepo.FindAll()
}

func (s *assetService) GetAssetByID(id string) (*models.Asset, error) {
	return s.assetRepo.FindByID(id)
}

func (s *assetService) CreateAsset(asset *models.Asset, userID uint) error {
	if !asset.Status.IsValid() {
		return errors.New("invalid status. Allowed: OK, Dismantle, Defect")
	}
	asset.AddedBy = userID
	return s.assetRepo.Create(asset)
}

func (s *assetService) UpdateAsset(id string, input models.Asset) (*models.Asset, error) {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update hanya field yang diizinkan
	asset.Name = input.Name
	asset.Type = input.Type
	asset.DeliveryDate = input.DeliveryDate
	asset.Status = input.Status
	asset.Location = input.Location
	asset.SerialNumber = input.SerialNumber

	err = s.assetRepo.Update(asset)
	return asset, err
}

func (s *assetService) DeleteAsset(id string, userID uint) error {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.assetRepo.SoftDelete(asset, userID)
}
