package services

import (
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/repositories"

	"github.com/gofiber/fiber/v2"
)

type JwtAuthService interface {
	CreateAccessToken(name string, email string, userId string, image string, emailVerified bool) (string, error)
	CreateRefreshToken(userId string) (string, error)
	RefreshAccessToken(refreshToken string) (string, error)
	VerifyJWT(tokenString string) (*entities.CustomClaim, error)
	HandleSetAccessTokenInCookie(c *fiber.Ctx, token string) error
	HandleSetRefreshTokenInCookie(c *fiber.Ctx, refreshToken string) error
	ClearAuthCookies(c *fiber.Ctx) error
	SetUserRepository(userRepo repositories.UserRepository)
	SetRefreshTokenRepository(refreshTokenRepo repositories.RefreshTokenRepository)
	InvalidateRefreshToken(token string) error
	InvalidateAllRefreshTokensForUser(userId string) error
}
