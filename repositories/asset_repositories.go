package repositories

import (
	"asset-maintenance/models"
	"time"

	"gorm.io/gorm"
)

type AssetRepository interface {
	FindAll() ([]models.Asset, error)
	FindByID(id string) (*models.Asset, error)
	Create(asset *models.Asset) error
	Update(asset *models.Asset) error
	SoftDelete(asset *models.Asset, deletedBy uint) error
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) FindAll() ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.Find(&assets).Error
	return assets, err
}

func (r *assetRepository) FindByID(id string) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.First(&asset, id).Error
	return &asset, err
}

func (r *assetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

func (r *assetRepository) Update(asset *models.Asset) error {
	return r.db.Save(asset).Error
}

func (r *assetRepository) SoftDelete(asset *models.Asset, deletedBy uint) error {
	now := time.Now()
	asset.DeletedBy = &deletedBy
	asset.DeletedAt = &now
	return r.db.Save(asset).Delete(asset).Error
}
