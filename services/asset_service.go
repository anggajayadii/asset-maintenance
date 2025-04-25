package services

import (
	"asset-maintenance/models"
	"asset-maintenance/repositories"
	"encoding/json"
	"errors"
	"log"
)

type AssetService interface {
	GetAllAssets() ([]models.Asset, error)
	GetAssetByID(id string) (*models.Asset, error)
	CreateAsset(asset *models.Asset, userID uint) error
	UpdateAsset(id string, input models.Asset, updatedBy uint) (*models.Asset, error) // Tambah updatedBy
	DeleteAsset(id string, userID uint) error
}

type assetService struct {
	assetRepo   repositories.AssetRepository
	historyRepo repositories.AssetHistoryRepository // Tambah historyRepo
}

// Revisi constructor untuk include historyRepo
func NewAssetService(
	assetRepo repositories.AssetRepository,
	historyRepo repositories.AssetHistoryRepository,
) AssetService {
	return &assetService{
		assetRepo:   assetRepo,
		historyRepo: historyRepo,
	}
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

	// Buat history setelah create
	if err := s.assetRepo.Create(asset); err != nil {
		return err
	}

	historyRecord := models.AssetHistory{
		AssetID:   asset.ID,
		Action:    "CREATE",
		ChangedBy: userID,
		NewData:   s.assetToJSON(asset),
	}

	return s.historyRepo.Create(historyRecord)
}

func (s *assetService) UpdateAsset(id string, input models.Asset, updatedBy uint) (*models.Asset, error) {
	// Dapatkan data lama sebelum update
	oldAsset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update field yang diizinkan
	updatedAsset := *oldAsset
	updatedAsset.Name = input.Name
	updatedAsset.Type = input.Type
	updatedAsset.DeliveryDate = input.DeliveryDate
	updatedAsset.Status = input.Status
	updatedAsset.Location = input.Location
	updatedAsset.SerialNumber = input.SerialNumber

	// Simpan perubahan
	if err := s.assetRepo.Update(&updatedAsset); err != nil {
		return nil, err
	}

	// Buat history record
	historyRecord := models.AssetHistory{
		AssetID:   updatedAsset.ID,
		Action:    "UPDATE",
		ChangedBy: updatedBy,
		OldData:   s.assetToJSON(oldAsset),
		NewData:   s.assetToJSON(&updatedAsset),
	}

	if err := s.historyRepo.Create(historyRecord); err != nil {
		log.Printf("Failed to create history record: %v", err)
		// Tidak return error karena update asset sudah berhasil
	}

	return &updatedAsset, nil
}

func (s *assetService) DeleteAsset(id string, userID uint) error {
	asset, err := s.assetRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Buat history sebelum delete
	historyRecord := models.AssetHistory{
		AssetID:   asset.ID,
		Action:    "DELETE",
		ChangedBy: userID,
		OldData:   s.assetToJSON(asset),
	}

	if err := s.historyRepo.Create(historyRecord); err != nil {
		log.Printf("Failed to create delete history: %v", err)
	}

	return s.assetRepo.SoftDelete(asset, userID)
}

// Helper untuk konversi asset ke JSON
func (s *assetService) assetToJSON(asset *models.Asset) string {
	jsonData, err := json.Marshal(asset)
	if err != nil {
		log.Printf("Error marshaling asset to JSON: %v", err)
		return ""
	}
	return string(jsonData)
}
