package models

import (
	"errors"

	"gorm.io/gorm"
)

type MaintenanceType string
type MaintenanceStatus string

const (
	// Maintenance Types
	Troubleshoot          MaintenanceType = "troubleshoot"
	PreventiveMaintenance MaintenanceType = "preventive_maintenance"

	// Maintenance Statuses
	StatusDone           MaintenanceStatus = "Done"
	StatusPending        MaintenanceStatus = "Pending"
	StatusRequiresAction MaintenanceStatus = "Requires Action"
)

type Maintenance struct {
	MaintID        uint              `gorm:"primaryKey" json:"maint_id"`
	AssetID        uint              `json:"asset_id"`
	Description    string            `json:"description"`
	UserID         uint              `json:"user_id"`
	Status         MaintenanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	MaintType      MaintenanceType   `gorm:"type:enum('troubleshoot','preventive_maintenance');not null" json:"maint_type"`
	Date           string            `json:"date"`            // ubah dari time.Time
	CompletionDate string            `json:"completion_date"` // ubah dari *time.Time

	// Relations
	Asset Asset `gorm:"foreignKey:AssetID"`
	User  User  `gorm:"foreignKey:UserID"`
}

// Validasi enum sebelum simpan
func (m *Maintenance) BeforeSave(tx *gorm.DB) (err error) {
	// Validasi MaintType
	switch m.MaintType {
	case Troubleshoot, PreventiveMaintenance:
	default:
		return errors.New("invalid maintenance type: must be 'troubleshoot' or 'preventive_maintenance'")
	}

	// Validasi Status
	switch m.Status {
	case StatusDone, StatusPending, StatusRequiresAction:
	default:
		return errors.New("invalid status: must be 'Done', 'Pending', or 'Requires Action'")
	}

	return nil
}
