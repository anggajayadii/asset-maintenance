package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type AssetStatus string

const (
	StatusOK        AssetStatus = "OK"
	StatusDismantle AssetStatus = "Dismantle"
	StatusDefect    AssetStatus = "Defect"
)

func (s AssetStatus) IsValid() bool {
	switch s {
	case StatusOK, StatusDismantle, StatusDefect:
		return true
	default:
		return false
	}
}

type DateOnly struct {
	time.Time
}

const customDateLayout = "2006-01-02"

// JSON Unmarshal
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1] // remove quotes
	t, err := time.Parse(customDateLayout, str)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

// JSON Marshal
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(customDateLayout))
}

// DB Scan (for SELECT)
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	default:
		return fmt.Errorf("cannot scan value %v into DateOnly", value)
	}
}

// DB Value (for INSERT/UPDATE)
func (d DateOnly) Value() (driver.Value, error) {
	return d.Time, nil
}

type Asset struct {
	AssetID      uint        `gorm:"primaryKey" json:"asset_id"`
	Name         string      `gorm:"not null" json:"name"`
	Type         string      `json:"type"`
	DeliveryDate *DateOnly   `json:"delivery_date"`
	Status       AssetStatus `json:"status"`
	Location     string      `json:"location"`
	SerialNumber string      `json:"serial_number"`

	AddedBy uint `json:"added_by"`
	// Tambahkan di model Asset:
	DeletedBy *uint      `json:"deleted_by" gorm:"default:null"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;default:null"`

	// Tambahkan relasi:
	DeletedByUser *User `gorm:"foreignKey:DeletedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Relations
	Maintenances []Maintenance  `gorm:"foreignKey:AssetID"`
	Histories    []AssetHistory `gorm:"foreignKey:AssetID"`
}
