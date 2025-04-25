package services

import (
	"asset-maintenance/models"
	"asset-maintenance/repositories"
)

type MaintenanceService interface {
	GetAllRecords() ([]models.Maintenance, error)
	GetRecordByID(id string) (*models.Maintenance, error)
	CreateRecord(maint *models.Maintenance) error
}

type maintenanceService struct {
	maintRepo repositories.MaintenanceRepository
}

func NewMaintenanceService(maintRepo repositories.MaintenanceRepository) MaintenanceService {
	return &maintenanceService{maintRepo: maintRepo}
}

func (s *maintenanceService) GetAllRecords() ([]models.Maintenance, error) {
	return s.maintRepo.FindAll()
}

func (s *maintenanceService) GetRecordByID(id string) (*models.Maintenance, error) {
	return s.maintRepo.FindByID(id)
}

func (s *maintenanceService) CreateRecord(maint *models.Maintenance) error {
	return s.maintRepo.Create(maint)
}
