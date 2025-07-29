package repository

import (
	"cardgame/domain"
	"errors"
	"time"

	"gorm.io/gorm"
)

// All 'token' parameters in this repository are expected to be the SHA-256 hash of the refresh token, not the raw token value.

type SQLRefreshTokenRepository struct {
	db *gorm.DB
}

func NewSQLRefreshTokenRepository(db *gorm.DB) *SQLRefreshTokenRepository {
	return &SQLRefreshTokenRepository{db: db}
}

func (r *SQLRefreshTokenRepository) CreateRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *SQLRefreshTokenRepository) GetRefreshTokenByToken(token string) (*domain.RefreshToken, error) {
	// token is now the hash
	var refreshToken domain.RefreshToken
	if err := r.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *SQLRefreshTokenRepository) GetRefreshTokensByUserID(userID string) ([]*domain.RefreshToken, error) {
	var refreshTokens []*domain.RefreshToken

	if err := r.db.Where("user_id = ?", userID).Find(&refreshTokens).Error; err != nil {
		return nil, err
	}

	return refreshTokens, nil
}

func (r *SQLRefreshTokenRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}

func (r *SQLRefreshTokenRepository) DeleteRefreshTokensByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.RefreshToken{}).Error
}

func (r *SQLRefreshTokenRepository) IsRefreshTokenValid(token string) (bool, error) {
	var refreshToken domain.RefreshToken

	if err := r.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if time.Now().After(refreshToken.ExpiresAt) {
		return false, nil
	}
	return true, nil
}
