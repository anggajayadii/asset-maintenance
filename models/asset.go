package models

import (
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

// type DateOnly struct {
// 	time.Time
// }

// const customDateLayout = "2006-01-02"

// func (d *DateOnly) UnmarshalJSON(b []byte) error {
// 	str := string(b)
// 	str = str[1 : len(str)-1] // remove quotes
// 	t, err := time.Parse(customDateLayout, str)
// 	if err != nil {
// 		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
// 	}
// 	d.Time = t
// 	return nil
// }

// func (d DateOnly) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(d.Time.Format(customDateLayout))
// }

// func (d *DateOnly) Scan(value interface{}) error {
// 	if value == nil {
// 		d.Time = time.Time{}
// 		return nil
// 	}
// 	switch v := value.(type) {
// 	case time.Time:
// 		d.Time = v
// 		return nil
// 	default:
// 		return fmt.Errorf("cannot scan value %v into DateOnly", value)
// 	}
// }

// func (d DateOnly) Value() (driver.Value, error) {
// 	return d.Time, nil
// }

type Asset struct {
	AssetID      uint        `gorm:"primaryKey" json:"asset_id"`
	Name         string      `gorm:"size:255;not null" json:"name"`
	Type         string      `gorm:"size:100" json:"type"`
	DeliveryDate string      `json:"delivery_date"`
	Status       AssetStatus `gorm:"size:20" json:"status"`
	Location     string      `gorm:"size:255" json:"location"`
	SerialNumber string      `gorm:"size:100;uniqueIndex" json:"serial_number"`

	// Audit fields
	AddedBy   uint       `json:"added_by"`
	AddedAt   time.Time  `gorm:"autoCreateTime" json:"added_at"`
	UpdatedBy *uint      `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedBy *uint      `json:"deleted_by"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	// Relations
	AddedByUser   *User          `gorm:"foreignKey:AddedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	UpdatedByUser *User          `gorm:"foreignKey:UpdatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	DeletedByUser *User          `gorm:"foreignKey:DeletedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Maintenances  []Maintenance  `gorm:"foreignKey:AssetID" json:"maintenances,omitempty"`
	Histories     []AssetHistory `gorm:"foreignKey:AssetID;constraint:OnDelete:CASCADE;" json:"histories,omitempty"`
}

// // BeforeCreate hook
// func (a *Asset) BeforeCreate(tx *gorm.DB) error {
// 	if !a.Status.IsValid() {
// 		return fmt.Errorf("invalid asset status: %s", a.Status)
// 	}

// 	if _, err := time.Parse("2006-01-02", a.DeliveryDate); err != nil {
// 		return fmt.Errorf("format delivery_date harus YYYY-MM-DD")
// 	}

// 	return nil
// }

// // BeforeUpdate hook
// func (a *Asset) BeforeUpdate(tx *gorm.DB) error {
// 	// Get current user ID from context
// 	if userID, ok := tx.Statement.Context.Value("user_id").(uint); ok {
// 		a.UpdatedBy = &userID
// 	}
// 	now := time.Now()
// 	a.UpdatedAt = &now
// 	return nil
// }

// // BeforeDelete hook (for soft delete)
// func (a *Asset) BeforeDelete(tx *gorm.DB) error {
// 	// Get current user ID from context
// 	if userID, ok := tx.Statement.Context.Value("user_id").(uint); ok {
// 		a.DeletedBy = &userID
// 	}
// 	now := time.Now()
// 	a.DeletedAt = &now

// 	// Save before delete for soft delete
// 	return tx.Save(a).Error
// }

// // AfterUpdate hook for history tracking
// func (a *Asset) AfterUpdate(tx *gorm.DB) error {
// 	if tx.Statement.Changed() {
// 		oldAsset, ok := tx.Statement.Context.Value("old_asset").(*Asset)
// 		if !ok {
// 			return nil
// 		}

// 		userID, _ := tx.Statement.Context.Value("user_id").(uint)

// 		historyService := tx.Statement.Context.Value("history_service").(services.AssetHistoryService)
// 		if historyService == nil {
// 			return nil
// 		}

// 		return historyService.RecordAssetHistory(
// 			a.AssetID,
// 			userID,
// 			"UPDATE",
// 			oldAsset,
// 			a,
// 			"Asset updated",
// 		)
// 	}
// 	return nil
// }

// // AfterCreate hook for history tracking
// func (a *Asset) AfterCreate(tx *gorm.DB) error {
// 	userID, _ := tx.Statement.Context.Value("user_id").(uint)

// 	historyService := tx.Statement.Context.Value("history_service").(services.AssetHistoryService)
// 	if historyService == nil {
// 		return nil
// 	}

// 	return historyService.RecordAssetHistory(
// 		a.AssetID,
// 		userID,
// 		"CREATE",
// 		nil,
// 		a,
// 		"Asset created",
// 	)
// }

// // AfterDelete hook for history tracking (soft delete)
// func (a *Asset) AfterDelete(tx *gorm.DB) error {
// 	userID, _ := tx.Statement.Context.Value("user_id").(uint)

// 	historyService := tx.Statement.Context.Value("history_service").(services.AssetHistoryService)
// 	if historyService == nil {
// 		return nil
// 	}

// 	return historyService.RecordAssetHistory(
// 		a.AssetID,
// 		userID,
// 		"DELETE",
// 		a,
// 		nil,
// 		"Asset deleted",
// 	)
// }
