package repositories

import (
	"asset-maintenance/models"

	"gorm.io/gorm"
)

type MaintenanceRepository interface {
	FindAll() ([]models.Maintenance, error)
	FindByID(id string) (*models.Maintenance, error)
	Create(maint *models.Maintenance) error
}

type maintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) MaintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) FindAll() ([]models.Maintenance, error) {
	var records []models.Maintenance
	err := r.db.Preload("Asset").Preload("User").Find(&records).Error
	return records, err
}

func (r *maintenanceRepository) FindByID(id string) (*models.Maintenance, error) {
	var maint models.Maintenance
	err := r.db.Preload("Asset").Preload("User").First(&maint, id).Error
	return &maint, err
}

func (r *maintenanceRepository) Create(maint *models.Maintenance) error {
	return r.db.Create(maint).Error
}
