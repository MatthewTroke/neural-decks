package repositories

import "cardgame/internal/domain/entities"

type RefreshTokenRepository interface {
	CreateRefreshToken(token *entities.RefreshToken) error
	GetRefreshTokenByToken(token string) (*entities.RefreshToken, error)
	GetRefreshTokensByUserID(userID string) ([]*entities.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokensByUserID(userID string) error
	IsRefreshTokenValid(token string) (bool, error)
}
