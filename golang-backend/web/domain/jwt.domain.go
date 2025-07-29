package domain

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct {
	Name                 string `json:"name"`           // Custom field for email
	Email                string `json:"email"`          // Custom field for email
	UserID               string `json:"user_id"`        // Custom field for user ID
	Image                string `json:"image"`          // Custom field for image
	EmailVerified        bool   `json:"email_verified"` // Custom field for email verification status
	jwt.RegisteredClaims        // Embeds the standard registered claims
}

type JwtAuthService interface {
	CreateAccessToken(name string, email string, userId string, image string, emailVerified bool) (string, error)
	CreateRefreshToken(userId string) (string, error)
	RefreshAccessToken(refreshToken string) (string, error)
	VerifyJWT(tokenString string) (*CustomClaim, error)
	HandleSetAccessTokenInCookie(c *fiber.Ctx, token string) error
	HandleSetRefreshTokenInCookie(c *fiber.Ctx, refreshToken string) error
	ClearAuthCookies(c *fiber.Ctx) error
	SetUserRepository(userRepo UserRepository)
	SetRefreshTokenRepository(refreshTokenRepo RefreshTokenRepository)
	InvalidateRefreshToken(token string) error
	InvalidateAllRefreshTokensForUser(userId string) error
}
