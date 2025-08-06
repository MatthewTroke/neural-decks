package interfaces

import "cardgame/domain/entities"

// AuthServiceInterface defines the contract for authentication services
type AuthServiceInterface interface {
	CreateAccessToken(name, email, userID, image string, emailVerified bool) (string, error)
	CreateRefreshToken(userID string) (string, error)
	InvalidateRefreshToken(token string) error
	InvalidateAllRefreshTokensForUser(userID string) error
	VerifyJWT(tokenString string) (*entities.CustomClaim, error)
}
