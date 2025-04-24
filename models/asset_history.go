package models

import "time"

type AssetHistory struct {
	HistoryID  uint      `gorm:"primaryKey" json:"history_id"`
	AssetID    uint      `json:"asset_id"`
	ChangeType string    `json:"change_type"`
	ChangedBy  uint      `json:"changed_by"`
	ChangedAt  time.Time `json:"changed_at" gorm:"autoCreateTime"`
	OldValue   string    `json:"old_value" gorm:"type:longtext"`
	NewValue   string    `json:"new_value" gorm:"type:longtext"`
	ChangeDesc string    `json:"change_desc" gorm:"type:longtext"`

	// Perbaikan relasi:
	Asset Asset `gorm:"foreignKey:AssetID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User  User  `gorm:"foreignKey:ChangedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
