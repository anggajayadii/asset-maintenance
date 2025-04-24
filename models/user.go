package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleEngineer Role = "engineer"
	RoleLogistik Role = "logistik"
	RoleManajer  Role = "manajer"
)

type User struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Role      Role      `gorm:"type:enum('engineer','logistik','manajer');not null" json:"role"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	AssetsAdded   []Asset        `gorm:"foreignKey:AddedBy"`
	AssetsDeleted []Asset        `gorm:"foreignKey:DeletedBy"`
	Maintenances  []Maintenance  `gorm:"foreignKey:UserID"`
	Changes       []AssetHistory `gorm:"foreignKey:ChangedBy"`
}

// Optional: Validasi role sebelum insert/update
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	switch u.Role {
	case RoleEngineer, RoleLogistik, RoleManajer:
		return nil
	default:
		return errors.New("invalid role: must be engineer, logistik, or manajer")
	}
}

// SetPassword meng-hash password dan menyimpannya ke user.Password
func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
