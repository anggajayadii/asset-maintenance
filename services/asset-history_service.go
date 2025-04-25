package services

import (
	"asset-maintenance/models"
	"asset-maintenance/repositories"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type AssetHistoryService interface {
	RecordAssetHistory(assetID uint, changedBy uint, changeType string, oldVal interface{}, newVal interface{}, desc string) error
	GetAssetHistories(assetID uint) ([]models.AssetHistory, error)
	GetHistoriesByDateRange(assetID uint, start, end time.Time) ([]models.AssetHistory, error)
	RecordAssetUpdate(asset *models.Asset, updatedAsset *models.Asset, changedBy uint) error // Diperbaiki dari Asset -> Asset
}

type assetHistoryService struct {
	historyRepo repositories.AssetHistoryRepository // Diperbaiki dari Asset -> Asset
}

func NewAssetHistoryService(historyRepo repositories.AssetHistoryRepository) AssetHistoryService { // Diperbaiki dari Asset -> Asset
	return &assetHistoryService{historyRepo: historyRepo}
}

func (s *assetHistoryService) RecordAssetHistory(assetID uint, changedBy uint, changeType string, oldVal interface{}, newVal interface{}, desc string) error {
	oldValue, _ := json.Marshal(oldVal)
	newValue, _ := json.Marshal(newVal)

	history := &models.AssetHistory{
		AssetID:    assetID,
		ChangeType: changeType,
		ChangedBy:  changedBy,
		OldValue:   string(oldValue),
		NewValue:   string(newValue),
		ChangeDesc: desc,
		ChangedAt:  time.Now(),
	}

	return s.historyRepo.CreateHistory(history)
}

func (s *assetHistoryService) GetAssetHistories(assetID uint) ([]models.AssetHistory, error) {
	return s.historyRepo.GetHistoriesByAssetID(assetID)
}

func (s *assetHistoryService) GetHistoriesByDateRange(assetID uint, start, end time.Time) ([]models.AssetHistory, error) {
	return s.historyRepo.GetHistoriesByDateRange(assetID, start, end)
}

func (s *assetHistoryService) RecordAssetUpdate(asset *models.Asset, updatedAsset *models.Asset, changedBy uint) error {
	// Use reflection to detect changes between original and updated asset
	originalVal := reflect.ValueOf(asset).Elem()
	updatedVal := reflect.ValueOf(updatedAsset).Elem()

	for i := 0; i < originalVal.NumField(); i++ {
		fieldName := originalVal.Type().Field(i).Name
		originalField := originalVal.Field(i).Interface()
		updatedField := updatedVal.Field(i).Interface()

		// Skip unexported fields and certain fields we don't want to track
		if fieldName == "Model" || fieldName == "CreatedAt" || fieldName == "DeletedAt" {
			continue
		}

		if !reflect.DeepEqual(originalField, updatedField) {
			desc := fmt.Sprintf("%s changed from %v to %v", fieldName, originalField, updatedField)
			err := s.RecordAssetHistory(
				asset.AssetID, // Sekarang bisa diakses karena tipe sudah benar
				changedBy,
				"UPDATE",
				originalField,
				updatedField,
				desc,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
