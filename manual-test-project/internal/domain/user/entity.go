package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents the domain entity for a user.
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
