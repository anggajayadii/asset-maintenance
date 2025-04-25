package repositories

import (
	"asset-maintenance/models"
	"time"

	"gorm.io/gorm"
)

type AssetHistoryRepository interface {
	CreateHistory(history *models.AssetHistory) error
	GetHistoriesByAssetID(assetID uint) ([]models.AssetHistory, error)
	GetHistoriesByDateRange(assetID uint, start, end time.Time) ([]models.AssetHistory, error)
}

type assetHistoryRepository struct {
	db *gorm.DB
}

func NewAssetHistoryRepository(db *gorm.DB) AssetHistoryRepository {
	return &assetHistoryRepository{db: db}
}

func (r *assetHistoryRepository) CreateHistory(history *models.AssetHistory) error {
	return r.db.Create(history).Error
}

func (r *assetHistoryRepository) GetHistoriesByAssetID(assetID uint) ([]models.AssetHistory, error) {
	var histories []models.AssetHistory
	err := r.db.Preload("User").Preload("Asset").
		Where("asset_id = ?", assetID).
		Order("changed_at DESC").
		Find(&histories).Error
	return histories, err
}

func (r *assetHistoryRepository) GetHistoriesByDateRange(assetID uint, start, end time.Time) ([]models.AssetHistory, error) {
	var histories []models.AssetHistory
	err := r.db.Preload("User").Preload("Asset").
		Where("asset_id = ? AND changed_at BETWEEN ? AND ?", assetID, start, end).
		Order("changed_at DESC").
		Find(&histories).Error
	return histories, err
}
