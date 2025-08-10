package entities

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Token     string         `gorm:"uniqueIndex;not null;type:text"`
	UserID    string         `gorm:"not null;type:varchar(255)"`
	ExpiresAt time.Time      `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
