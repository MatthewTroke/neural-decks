package domain

import (
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(token *RefreshToken) error
	GetRefreshTokenByToken(token string) (*RefreshToken, error)
	GetRefreshTokensByUserID(userID string) ([]*RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokensByUserID(userID string) error
	IsRefreshTokenValid(token string) (bool, error)
}

type RefreshToken struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Token     string         `gorm:"uniqueIndex;not null;type:text"`
	UserID    string         `gorm:"not null;type:varchar(255)"`
	ExpiresAt time.Time      `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
